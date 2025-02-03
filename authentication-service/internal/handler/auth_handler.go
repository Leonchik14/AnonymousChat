package handler

import (
	"context"
	"net/http"

	"authentication-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler - структура обработчика
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler - конструктор обработчика аутентификации
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register - обработчик регистрации пользователя
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	// Вызываем сервис для регистрации пользователя
	err := h.authService.RegisterUser(context.Background(), req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Пользователь успешно зарегистрирован"})
}

// Login - обработчик входа пользователя
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	// Вызываем сервис для авторизации
	accessToken, refreshToken, err := h.authService.LoginUser(context.Background(), req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Неверный email или пароль"})
	}

	// Отправляем токен в ответе
	return c.JSON(fiber.Map{"accessToken": accessToken, "refreshToken": refreshToken})
}

// Logout - обработчик выхода пользователя
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Получаем токен из заголовка Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Отсутствует токен"})
	}

	// Вызываем сервис для добавления токена в blacklist
	err := h.authService.LogoutUser(context.Background(), authHeader)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка выхода"})
	}

	return c.JSON(fiber.Map{"message": "Вы успешно вышли из системы"})
}

// Profile - обработчик получения профиля пользователя (требует JWT)
/*func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	// Извлекаем userID из JWT-токена
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Ошибка получения userID"})
	}

	return c.JSON(fiber.Map{"message": "Ваш профиль", "userID": userID})
}*/

// ForgotPassword - обработчик запроса на восстановление пароля
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	err := h.authService.ForgotPassword(context.Background(), req.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Ссылка для сброса пароля отправлена на email"})
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Отсутствует токен"})
	}

	userID, err := h.authService.ValidateToken(context.Background(), token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Недействительный токен"})
	}

	err = h.authService.VerifyUser(context.Background(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка подтверждения"})
	}

	return c.JSON(fiber.Map{"message": "Email подтвержден!"})
}

// SetupRoutes - регистрация маршрутов
func (h *AuthHandler) SetupRoutes(app *fiber.App) {
	app.Post("/auth/register", h.Register)
	app.Post("/auth/login", h.Login)
	app.Post("/auth/logout", h.Logout)
	app.Get("/auth/verify", h.VerifyEmail)
	app.Get("/auth/forgot-password", h.ForgotPassword)
}
