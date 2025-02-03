package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"matchmaking-service/internal/service"
)

// MatchmakingHandler - обработчик HTTP-запросов для поиска собеседников
type MatchmakingHandler struct {
	matchmakingService *service.MatchmakingService
}

// NewMatchmakingHandler - конструктор
func NewMatchmakingHandler(matchmakingService *service.MatchmakingService) *MatchmakingHandler {
	return &MatchmakingHandler{matchmakingService: matchmakingService}
}

// StartMatchmaking - обработчик запроса на поиск собеседника
func (h *MatchmakingHandler) StartMatchmaking(c *fiber.Ctx) error {
	// Получаем userID из запроса (например, из заголовка или токена)
	userIDStr := c.Get("X-User-ID") // Клиент должен передавать заголовок "X-User-ID"
	if userIDStr == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Необходимо передать X-User-ID"})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат X-User-ID"})
	}

	log.Printf("🔍 Пользователь %d начал поиск собеседника...", userID)

	// Ищем собеседника
	err = h.matchmakingService.FindMatch(context.Background(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Ответ клиенту: подтверждение запроса
	return c.JSON(fiber.Map{"message": "Поиск собеседника запущен, ожидайте."})
}
