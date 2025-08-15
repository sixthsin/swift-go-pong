package websocket

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"pong-api-v1/internal/game"

	"github.com/gorilla/websocket"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	playerId := rand.Text()
	newPlayer := &game.Player{
		Conn:  conn,
		Id:    playerId,
		Name:  "name",
		X:     0.5,
		Score: 0,
	}
	newGame := game.InitGame()
	newGame.AddPlayer(newPlayer) // Should be player id

	handlePlayer(conn, newGame, playerId)
}

func handlePlayer(conn *websocket.Conn, newGame *game.Game, playerId string) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var input game.PlayerInput
		if err := json.Unmarshal(msg, &input); err != nil {
			log.Printf("JSON unmarshall error:%v", err)
			continue
		}

		switch input.Type {
		case "racket_move":
			newGame.Players[playerId].X = input.X
		}
	}
}
