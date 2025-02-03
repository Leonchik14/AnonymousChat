package main

import (
	"authentication-service/internal/app"
)

func main() {
	application := app.NewApp() // 🔹 Создаем приложение
	application.Run()           // 🔹 Запускаем сервер
}
