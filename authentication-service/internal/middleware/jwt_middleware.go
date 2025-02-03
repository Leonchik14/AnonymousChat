package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"authentication-service/internal/service"
	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware - структура для middleware
type JWTMiddleware struct {
	authService *service.AuthService
}

// NewJWTMiddleware - конструктор middleware
func NewJWTMiddleware(authService *service.AuthService) *JWTMiddleware {
	return &JWTMiddleware{authService: authService}
}

// MiddlewareJWT - middleware для проверки JWT токена
func (m *JWTMiddleware) MiddlewareJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем токен из заголовка Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Отсутствует заголовок Authorization"})
		}

		// Заголовок должен быть вида "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Неверный формат токена"})
		}
		token := parts[1]

		// Валидация токена через сервис
		userID, err := m.authService.ValidateToken(context.Background(), token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Недействительный или истекший токен"})
		}

		// Добавляем userID в локальный контекст запроса
		c.Locals("userID", userID)
		return c.Next()
	}
}

// ExtractUserID - извлекает userID из локального контекста запроса
func ExtractUserID(c *fiber.Ctx) (int64, error) {
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return 0, errors.New("не удалось получить userID из контекста")
	}
	return userID, nil
}
