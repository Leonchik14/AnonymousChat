package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"chat-service/internal/repository"
	"chat-service/pkg/models"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	chatRepo *repository.ChatRepository
	clients  map[int64]map[*websocket.Conn]bool // WebSocket-клиенты (chat_id -> conn)
}

func NewChatHandler(chatRepo *repository.ChatRepository) *ChatHandler {
	return &ChatHandler{
		chatRepo: chatRepo,
		clients:  make(map[int64]map[*websocket.Conn]bool),
	}
}

func (h *ChatHandler) WebSocketHandler(c *websocket.Conn) {
	defer c.Close()

	// Получаем chat_id из пути
	chatIDStr := c.Params("chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Println("❌ Ошибка: Неверный chat_id")
		return
	}

	// Регистрируем нового WS-клиента
	if h.clients[chatID] == nil {
		h.clients[chatID] = make(map[*websocket.Conn]bool)
	}
	h.clients[chatID][c] = true
	log.Printf("✅ Пользователь подключился к чату %d", chatID)

	for {
		// Читаем сырое сообщение
		_, raw, err := c.ReadMessage()
		if err != nil {
			log.Printf("❌ Отключение клиента от чата %d: %v", chatID, err)
			delete(h.clients[chatID], c)
			break
		}

		// Парсим JSON из фронтенда
		var wsMsg models.WsMessage
		if err := json.Unmarshal(raw, &wsMsg); err != nil {
			log.Println("❌ Неверный формат WS-сообщения:", err)
			continue
		}

		// Проставляем серверное время
		wsMsg.Timestamp = time.Now().UTC()

		// Сохраняем в БД
		modelMsg := models.Message{
			ChatID:    chatID,
			SenderID:  wsMsg.SenderID,
			Content:   wsMsg.Content,
			CreatedAt: wsMsg.Timestamp,
		}
		if err := h.chatRepo.SaveMessage(context.Background(), &modelMsg); err != nil {
			log.Printf("❌ Ошибка сохранения сообщения: %v", err)
			continue
		}

		// Рассылаем всем подключённым клиентам
		broadcast, _ := json.Marshal(wsMsg)
		for client := range h.clients[chatID] {
			if err := client.WriteMessage(websocket.TextMessage, broadcast); err != nil {
				log.Println("❌ Ошибка отправки сообщения:", err)
			}
		}
	}
}

func (h *ChatHandler) GetChatHistory(c *fiber.Ctx) error {
	chatID, err := strconv.ParseInt(c.Params("chat_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный chat_id"})
	}

	messages, err := h.chatRepo.GetChatHistory(context.Background(), chatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка загрузки истории чата"})
	}

	return c.JSON(messages)
}
