package main

import "matchmaking-service/internal/app"

func main() {
	application := app.NewApp() // 🔹 Создаем приложение
	application.Run()           // 🔹 Запускаем сервер
}
