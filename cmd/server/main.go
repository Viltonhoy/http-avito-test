package main

import (
	"context"
	"http-avito-test/internal/exchanger"
	"http-avito-test/internal/server"
	"http-avito-test/internal/storage"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zap.NewDevelopment: %v", err)
	}
	defer logger.Sync()

	if err := godotenv.Load("../../.env"); err != nil {
		logger.Debug("No .env file found", zap.Error(err))
	}

	ctx := context.Background()

	storage, err := storage.NewStorage(ctx, logger)
	if err != nil {
		logger.Fatal("failed to create storage instance", zap.Error(err))
	}

	e := exchanger.New()

	srv, err := server.New(
		logger,
		storage,
		storage.Close,
		e,
	)

	if err != nil {
		logger.Fatal("failed to create http server instance", zap.Error(err))
	}

	go func() {
		mux := http.NewServeMux()

		mux.Handle("/", http.FileServer(http.Dir("../../file_storage")))

		log.Println("Запуск сервера на http://localhost:4000")
		err := http.ListenAndServe(":4000", mux)
		log.Fatal(err)
	}()

	err = srv.Start()
	if err != nil {
		logger.Fatal("failed to start or shutdown server", zap.Error(err))
	}
}
