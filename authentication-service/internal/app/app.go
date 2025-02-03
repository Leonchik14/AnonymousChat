package app

import (
	"log"
	"os"
	"time"

	"authentication-service/internal/handler"
	"authentication-service/internal/repository"
	"authentication-service/internal/service"
	"authentication-service/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// App - структура приложения
type App struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	AuthService *service.AuthService
	FiberApp    *fiber.App
}

// NewApp - инициализация приложения
func NewApp() *App {
	// 🔹 Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env файл не найден, используем системные переменные")
	}

	// 🔹 Читаем переменные окружения
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	// 🔹 Формируем строку подключения к MySQL
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"

	// 🔹 Подключаем базу данных через GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к MySQL: %v", err)
	}

	// 🔹 Автоматически создаем таблицы (если их нет)
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("❌ Ошибка миграции базы данных: %v", err)
	}
	log.Println("✅ Таблицы созданы или уже существуют")

	// 🔹 Подключаем Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	// 🔹 Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)

	// 🔹 Создаем сервисы
	emailService := service.NewEmailService()
	authService := service.NewAuthService(userRepo, redisRepo, emailService, jwtSecret, 15*time.Minute, 7*24*time.Hour)

	// 🔹 Инициализируем Fiber (HTTP-сервер)
	app := fiber.New()

	// 🔹 Создаем обработчики (Handlers)
	authHandler := handler.NewAuthHandler(authService)

	// 🔹 Настраиваем маршруты
	authHandler.SetupRoutes(app)

	return &App{
		DB:          db,
		RedisClient: redisClient,
		AuthService: authService,
		FiberApp:    app,
	}
}

// Run - запуск сервера
func (a *App) Run() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Сервер запущен на порту %s", port)
	log.Fatal(a.FiberApp.Listen(":" + port))
}
