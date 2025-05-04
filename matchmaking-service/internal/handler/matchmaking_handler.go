package handler

import (
	"bufio"
	"context"
	"fmt"
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
	// 2) Запускаем поиск — пусть в сервисе это кладёт в очередь и
	// возвращает канал, куда придёт chatID, когда пара найдена.
	matchCh, err := h.matchmakingService.FindMatch(context.Background(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().Response.Header.Set("Transfer-Encoding", "chunked")
	c.Context().Response.Header.Set("X-Accel-Buffering", "no")

	c.Context().Response.SetBodyStreamWriter(func(w *bufio.Writer) {
		if chatID, ok := <-matchCh; ok {
			fmt.Fprintf(w, "event: match\ndata: %d\n\n", chatID)
			w.Flush()
		}
	})
	return nil
}
