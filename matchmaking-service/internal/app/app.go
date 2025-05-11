package app

import (
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"

	"matchmaking-service/internal/grpc/chatpb"
	"matchmaking-service/internal/handler"
	"matchmaking-service/internal/repository"
	"matchmaking-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

// App - структура для хранения компонентов приложения
type App struct {
	FiberApp *fiber.App
}

func NewApp() *App {
	// 🔹 Загружаем переменные окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env файл не найден, используем переменные окружения")
	}

	// 🔹 Читаем переменные окружения
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	chatServiceHost := os.Getenv("CHAT_SERVICE_HOST")
	chatServicePort := os.Getenv("CHAT_SERVICE_PORT")
	// 🔹 Подключаем Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})
	redisRepo := repository.NewRedisRepository(redisClient)

	chatConn, err := grpc.NewClient(chatServiceHost+":"+chatServicePort, grpc.WithTransportCredentials(insecure.NewCredentials())) // gRPC клиент
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к chat-service: %v", err)
	}
	chatClient := chatpb.NewChatServiceClient(chatConn) // Создаем gRPC-клиент

	// 🔹 Создаем сервис поиска собеседников
	matchmakingService := service.NewMatchmakingService(redisRepo, chatClient)

	// 🔹 Инициализируем HTTP-сервер
	app := fiber.New()
	matchmakingHandler := handler.NewMatchmakingHandler(matchmakingService)

	// 🔹 Регистрация маршрутов
	app.Get("api/matchmaking/start", matchmakingHandler.StartMatchmaking)

	return &App{
		FiberApp: app,
	}
}

// Run - запуск HTTP сервера
func (a *App) Run() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 Matchmaking Service запущен на порту %s", port)
	log.Fatal(a.FiberApp.Listen(":" + port))
}
