package main

import (
	generatetable "http-avito-test/internal/generateTable"
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

	var s, _ = generatetable.NewStore(logger)

	generatetable.AddGeneratedTable(s, 500, 1000000)
}
