package game

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	PlayerStartScore    = 0
	PlayerDefaultX      = 0.5
	ballDefaultX        = 0.5
	ballDefaultY        = 0.5
	ballDefaultDX       = 0.004
	ballDefaultDY       = 0.002
	gameMaxPlayersCount = 2
)

func InitGame() *Game {
	return &Game{
		Id:      generateGameId(),
		Players: make(map[string]*Player),
		Ball: Ball{
			X:  ballDefaultX,
			Y:  ballDefaultY,
			DX: ballDefaultDX,
			DY: ballDefaultDY,
		},
		Started: false,
		Mu:      &sync.Mutex{},
	}
}

func (g *Game) AddPlayer(p *Player) error {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	if len(g.Players) <= gameMaxPlayersCount {
		g.Players[p.Id] = p
		return nil
	}
	p.Conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Game is full"}`))
	return errors.New("Game is full")
}

func generateGameId() string {
	return uuid.New().String()
}

func (g *Game) Run() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		g.Update()
		g.BroadcastGameState()
	}
}

func (g *Game) Update() {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	g.Ball.X += g.Ball.DX
	g.Ball.Y += g.Ball.DY

	if g.Ball.X < 0 || g.Ball.X > 1 {
		g.Ball.DX *= -1
	}
	if g.Ball.Y < 0 || g.Ball.Y > 1 {
		g.Ball.DY *= -1
	}
}

func (g *Game) BroadcastGameState() {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	state, _ := json.Marshal(g)
	for _, player := range g.Players {
		if err := player.Conn.WriteMessage(websocket.TextMessage, state); err != nil {
			log.Printf("JSON error: %v", err.Error())
		}
	}
}
