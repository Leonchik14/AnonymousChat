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

	chatIDStr := c.Params("chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Println("❌ Ошибка: Неверный chat_id")
		return
	}

	if h.clients[chatID] == nil {
		h.clients[chatID] = make(map[*websocket.Conn]bool)
	}
	h.clients[chatID][c] = true
	log.Printf("✅ Пользователь подключился к чату %d", chatID)

	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			log.Printf("❌ Отключение клиента от чата %d: %v", chatID, err)
			delete(h.clients[chatID], c)
			break
		}

		var wsMsg models.WsMessage
		if err := json.Unmarshal(raw, &wsMsg); err != nil {
			log.Println("❌ Неверный формат WS-сообщения:", err)
			continue
		}

		wsMsg.Timestamp = time.Now().UTC()

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

		for client := range h.clients[chatID] {
			if err := client.WriteJSON(modelMsg); err != nil {
				log.Println("❌ Ошибка отправки JSON-сообщения:", err)
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

func (h *ChatHandler) GetAllChats(c *fiber.Ctx) error {
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Необходимо передать X-User-ID"})
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат X-User-ID"})
	}

	chats, err := h.chatRepo.GetUserChats(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(chats)
}
