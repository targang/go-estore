package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go_store/config"
	"go_store/internal/app"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg, err := config.New()

	if err != nil {
		log.Fatalf("can not get application config: %s", err)
	}

	var logger *zap.Logger

	logger, err = zap.NewProduction()

	if err != nil {
		log.Fatalf("can not initialize logger: %s", err)
	}

	app.Run(logger, cfg)
}
