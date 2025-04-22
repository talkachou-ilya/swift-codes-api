package main

import (
	"swift-codes-api/internal/app"
	"swift-codes-api/internal/config"
)

func main() {
	cfg := config.Load()
	application := app.New(cfg)
	app.Start(application)
}
