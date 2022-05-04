package main

import (
	"http-avito-test/internal/server"
	"http-avito-test/internal/storage"
	"log"
	"net/http"

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
	sugar := logger.Sugar()

	var s, _ = storage.NewStore(sugar)
	h := server.Handler{
		Store: s,
	}
	http.HandleFunc("/read", h.ReadUser)
	http.HandleFunc("/update", h.AccountFunding)
	http.HandleFunc("/transf", h.TransferCommand)

	port := ":9090"
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListernAndServe", err)
	}
}
