package server

import (
	"encoding/json"
	"log"
	"net/http"
	"pong-api-v1/internal/game"
	"pong-api-v1/internal/storage"

	"github.com/gorilla/websocket"
)

type Handlers struct {
	Storage  *storage.MemoryStorage
	Upgrader *websocket.Upgrader
}

func NewHandlers(storage *storage.MemoryStorage) *Handlers {
	return &Handlers{
		Storage: storage,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handlers) CreateGame(w http.ResponseWriter, r *http.Request) {
	if enableCORS(&w, r) {
		return
	}
	gameId := r.URL.Query().Get("game_id")
	if gameId == "" {
		log.Println("game_id is empty")
		http.Error(w, "game_id is required", http.StatusBadRequest)
		return
	}

	if _, err := h.Storage.Get(gameId); err == nil {
		log.Printf("game %s already exists", gameId)
		http.Error(w, "game already exists", http.StatusConflict)
		return
	}

	game := game.InitGame(gameId)

	h.Storage.Save(game)
	log.Printf("created game %s", gameId)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"game created"}`))
}

func (h *Handlers) JoinGame(w http.ResponseWriter, r *http.Request) {
	if enableCORS(&w, r) {
		return
	}
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	gameID := r.URL.Query().Get("game_id")
	if gameID == "" {
		http.Error(w, `{"error": "game_id is required"}`, http.StatusBadRequest)
		return
	}

	game, err := h.Storage.Get(gameID)
	if err != nil {
		http.Error(w, `{"error": "game not found"}`, http.StatusNotFound)
		return
	}

	if len(game.Players) >= 2 {
		http.Error(w, `{"error": "game is full"}`, http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "can join",
		"game_id": gameID,
	})
}

func enableCORS(w *http.ResponseWriter, r *http.Request) bool {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return true
	}
	return false
}

// func generateGameId() string {
// 	return uuid.New().String()
// }
