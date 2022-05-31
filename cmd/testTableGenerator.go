package main

import (
	"context"
	"http-avito-test/internal/storage"
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zap.NewDevelopment: %v", err)
	}
	defer logger.Sync()

	var s, _ = storage.NewStore(context.Background(), logger)

	storage.AddGeneratedTable(s, 5, 100)
}
