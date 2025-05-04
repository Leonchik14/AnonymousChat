package main

import "matchmaking-service/internal/app"

func main() {
	application := app.NewApp()
	application.Run()
}
