package main

import (
	"chat-service/internal/app"
)

func main() {
	application := app.NewApp()
	application.Run()
}
