package main

import (
	"log"
	"pong-api-v1/internal/websocket"

	"net/http"
)

func main() {

	http.HandleFunc("/ws", websocket.HandleWebSocket)
	log.Println("Started server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
