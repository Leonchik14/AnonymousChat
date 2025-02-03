package handler

import (
	"context"
	"log"
	"strconv"

	"chat-service/internal/repository"
	"chat-service/pkg/models"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// ChatHandler - обработчик чатов
type ChatHandler struct {
	chatRepo *repository.ChatRepository
	clients  map[int64]map[*websocket.Conn]bool // WebSocket-клиенты (chat_id -> conn)
}

// NewChatHandler - конструктор обработчика
func NewChatHandler(chatRepo *repository.ChatRepository) *ChatHandler {
	return &ChatHandler{
		chatRepo: chatRepo,
		clients:  make(map[int64]map[*websocket.Conn]bool),
	}
}

// WebSocketHandler - обработка WebSocket соединений
func (h *ChatHandler) WebSocketHandler(c *websocket.Conn) {
	defer c.Close()

	// Получаем chat_id из параметров запроса
	chatIDStr := c.Params("chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Println("❌ Ошибка: Неверный chat_id")
		return
	}

	// Подключаем клиента к чату
	if h.clients[chatID] == nil {
		h.clients[chatID] = make(map[*websocket.Conn]bool)
	}
	h.clients[chatID][c] = true
	log.Printf("✅ Пользователь подключился к чату %d", chatID)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("❌ Отключение клиента от чата %d", chatID)
			delete(h.clients[chatID], c)
			break
		}

		message := models.Message{
			ChatID:  chatID,
			Content: string(msg),
		}
		if err := h.chatRepo.SaveMessage(context.Background(), &message); err != nil {
			log.Printf("❌ Ошибка сохранения сообщения: %v", err)
			continue
		}

		for client := range h.clients[chatID] {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("❌ Ошибка отправки сообщения:", err)
			}
		}
	}
}

// GetChatHistory - обработчик истории чата
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
