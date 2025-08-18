package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pong-api-v1/internal/game"

	"github.com/gorilla/websocket"
)

func (h *Handlers) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	gameId := r.URL.Query().Get("game_id")
	if gameId == "" {
		log.Println("game_id is empty")
		// err to client
		return
	}

	foundGame, err := h.Storage.Get(gameId)
	if err != nil {
		log.Printf("Game %s not found", gameId)
		// err to client
		return
	}

	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		//err to client
		return
	}
	// defer conn.Close()

	var playerId string
	switch len(foundGame.Players) {
	case 0:
		playerId = "player1"
	case 1:
		playerId = "player2"
	default:
		sendErrorMessage(conn, "game is full")
		log.Println("Connection error: game is full")
		conn.Close()
		return
	}

	if err := foundGame.AddPlayer(
		&game.Player{
			Conn:  conn,
			Id:    playerId,
			X:     game.PlayerDefaultX,
			Score: game.PlayerStartScore},
	); err != nil {
		sendErrorMessage(conn, err.Error())
		conn.Close()
		return
	}

	foundGame.WaitForPlayers()

	go foundGame.Run()
	go h.handlePlayer(conn, foundGame, playerId)
}

func (h *Handlers) handlePlayer(conn *websocket.Conn, currentGame *game.Game, playerId string) {
	defer func() {
		currentGame.Mu.Lock()
		defer currentGame.Mu.Unlock()

		if player, exists := currentGame.Players[playerId]; exists {
			player.Conn.Close()
			delete(currentGame.Players, playerId)
		}
	}()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			handlerConnectionError(err, playerId)
		}
		log.Printf("Received message from %s: %s", playerId, string(msg))

		if messageType == -1 { // websocket.ClosingMessage
			log.Printf("Player %s disconnected", playerId)
			return
		}

		var input game.PlayerInput
		if err := json.Unmarshal(msg, &input); err != nil {
			sendErrorMessage(conn, "invalid JSON format")
			log.Printf("JSON unmarshall error:%v", err)
			continue
		}

		if err := handleInput(conn, currentGame, playerId, input); err != nil {
			log.Printf("Movement error for %s: %v", playerId, err)
			continue
		}
	}
}

func handleInput(conn *websocket.Conn, currentGame *game.Game, playerId string, input game.PlayerInput) error {
	if input.Type == "racket_move" {
		if input.X < 0 || input.X > 1 {
			sendErrorMessage(conn, "position must be between 0 and 1")
			return fmt.Errorf("invalid position by %s: %f", playerId, input.X)
		}

		currentGame.Mu.Lock()
		defer currentGame.Mu.Unlock()

		if player, exists := currentGame.Players[playerId]; exists {
			player.X = input.X
		} else {
			sendErrorMessage(conn, "player not found")
			return fmt.Errorf("player %s not found in game", playerId)
		}
	}

	return nil
}

func handlerConnectionError(err error, playerId string) {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		log.Printf("Unexpected disconnect from %s: %v", playerId, err)
	} else {
		log.Printf("Player %s disconnected: %v", playerId, err)
	}
}

func sendErrorMessage(conn *websocket.Conn, msg string) {
	errMsg := map[string]string{"error": msg}
	if err := conn.WriteJSON(errMsg); err != nil {
		log.Printf("Failed to send error message: %v", err)
	}
}
