package main

import (
	"context"
	"http-avito-test/internal/storage"
	"log"
	"os"
	"strconv"

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
		return
	}

	values := os.Args[1:]
	userCount, _ := strconv.Atoi(values[0])
	totalRecordCount, _ := strconv.Atoi(values[1])

	s, err := storage.NewStorage(context.Background(), logger)
	if err != nil {
		return
	}

	AddGeneratedTableData(s, userCount, totalRecordCount)
}
