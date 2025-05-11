package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"matchmaking-service/internal/service"
)

type MatchmakingHandler struct {
	matchmakingService *service.MatchmakingService
}

func NewMatchmakingHandler(matchmakingService *service.MatchmakingService) *MatchmakingHandler {
	return &MatchmakingHandler{matchmakingService: matchmakingService}
}

// StartMatchmaking - обработчик запроса на поиск собеседника
func (h *MatchmakingHandler) StartMatchmaking(c *fiber.Ctx) error {
	// 1) Авторизация
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Необходимо передать X-User-ID"})
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат X-User-ID"})
	}

	log.Printf("🔍 Пользователь %d встал в очередь...", userID)

	matchCh, err := h.matchmakingService.FindMatch(context.Background(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 2) Ждём появления chatID в канале
	chatID, ok := <-matchCh
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "канал закрыт без совпадения"})
	}

	return c.JSON(fiber.Map{
		"event": "match",
		"data":  chatID,
	})
}
