package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"pong-api-v1/internal/game"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	currentGame     *game.Game
	gameInitializer sync.Once
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

	gameInitializer.Do(func() {
		currentGame = game.InitGame()
	})

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// defer conn.Close()

	playerId := uuid.New().String()
	newPlayer := &game.Player{
		Conn:  conn,
		Id:    playerId,
		X:     game.PlayerDefaultX,
		Score: game.PlayerStartScore,
	}

	if err := currentGame.AddPlayer(newPlayer); err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		conn.Close()
		return
	}

	for {
		if len(currentGame.Players) != 2 {
			continue
		} else {
			break
		}
	}

	go currentGame.Run()
	go handlePlayer(conn, currentGame, playerId)
}

func handlePlayer(conn *websocket.Conn, currentGame *game.Game, playerId string) {
	for {
		_, msg, err := conn.ReadMessage()
		log.Printf("Recieved request:%s", string(msg))
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
			if input.X < 0 || input.X > 1 {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Invalid position"}`))
				continue
			}
			currentGame.Mu.Lock()
			currentGame.Players[playerId].X = input.X
			currentGame.Mu.Unlock()
		}
	}
}
