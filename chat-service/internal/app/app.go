package app

import (
	"log"
	"os"

	"chat-service/internal/grpc"
	"chat-service/internal/handler"
	"chat-service/internal/repository"
	"chat-service/internal/service"
	"chat-service/pkg/models"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// App - структура приложения
type App struct {
	FiberApp *fiber.App
}

func NewApp() *App {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env файл не найден, используем переменные окружения")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к MySQL: %v", err)
	}

	if err := db.AutoMigrate(&models.Chat{}, &models.Message{}); err != nil {
		log.Fatalf("❌ Ошибка миграции базы данных: %v", err)
	}
	log.Println("✅ Таблицы созданы или уже существуют")

	// 🔹 Создаем репозитории
	chatRepo := repository.NewChatRepository(db)

	// 🔹 Создаем сервис
	chatService := service.NewChatService(chatRepo)

	// 🔹 Запускаем gRPC-сервер (асинхронно)
	go grpc.RunGRPCServer(chatService)

	// 🔹 Создаем HTTP-сервер с Fiber
	app := fiber.New()
	chatHandler := handler.NewChatHandler(chatRepo)

	app.Get("/ws/chat/:chat_id", websocket.New(chatHandler.WebSocketHandler))
	app.Get("/api/chat/history/:chat_id", chatHandler.GetChatHistory)

	return &App{
		FiberApp: app,
	}
}

// Run - запуск HTTP сервера
func (a *App) Run() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("🚀 Chat Service запущен на порту %s", port)
	log.Fatal(a.FiberApp.Listen(":" + port))
}
