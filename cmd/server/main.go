package main

import (
	"http-avito-test/internal/server"
	"log"
	"net/http"
)

func main() {
	h := server.Handler{}
	http.HandleFunc("/add", h.ReadUser)

	port := ":9090"
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListernAndServe", err)
	}

}
