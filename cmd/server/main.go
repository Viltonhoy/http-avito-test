package main

import (
	"context"
	"http-avito-test/internal/server"
	"http-avito-test/internal/storage"
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zap.NewDevelopment: %v", err)
	}
	defer logger.Sync()

	ctx := context.Background()

	storage, err := storage.NewStore(ctx, logger)
	if err != nil {
		logger.Fatal("failed to create storage instance", zap.Error(err))
	}

	srv, err := server.New(
		logger,
		storage.Close,
	)

	if err != nil {
		logger.Fatal("failed to create http server instance", zap.Error(err))
	}

	err = srv.Start()
	if err != nil {
		logger.Fatal("failed to start or shutdown server", zap.Error(err))
	}
}
