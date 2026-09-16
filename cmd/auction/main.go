package main

import (
	"auction-house/internal/auction"
	"auction-house/internal/httpapi"
	"auction-house/internal/storage"
	"log"
	"net/http"
)

func runServer() error {
	repo := storage.NewMemoryRepository()
	service := auction.NewService(repo)
	handler := httpapi.NewHandler(service)

	router := handler.Routes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("server listening on :8080")

	return server.ListenAndServe()
}

func main() {
	if err := runServer(); err != nil {
		log.Fatal(err)
	}
}
