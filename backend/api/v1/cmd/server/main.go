package main

import (
	"log"
	"pong-api-v1/internal/server"
	"pong-api-v1/internal/storage"

	"net/http"
)

func main() {
	storage := storage.NewStorage()
	handlers := server.NewHandlers(storage)

	http.HandleFunc("/ws", handlers.HandleWebSocket)
	http.HandleFunc("/api/game/create", handlers.CreateGame)
	http.HandleFunc("/api/game/join", handlers.JoinGame)

	log.Println("Started server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
