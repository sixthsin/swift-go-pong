package game

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	gameMaxPlayersCount = 2
	ticksPerSecond      = 60
	topRacketY          = 0.05
	bottomRacketY       = 0.95
	racketWidth         = 0.2
	racketHeight        = 0.01
)

var (
	PlayerStartScore = 0
	PlayerDefaultX   = 0.5
	ballDefaultX     = 0.5
	ballDefaultY     = 0.5
	ballDefaultDX    = 0.004
	ballDefaultDY    = 0.002
)

func InitGame(gameId string) *Game {
	return &Game{
		Id:      gameId,
		Players: make(map[string]*Player),
		Ball: Ball{
			X:  ballDefaultX,
			Y:  ballDefaultY,
			DX: ballDefaultDX,
			DY: ballDefaultDY,
		},
		Started:       false,
		Mu:            &sync.Mutex{},
		TopRacketY:    topRacketY,
		BottomRacketY: bottomRacketY,
	}
}

func (g *Game) AddPlayer(p *Player) error {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	if len(g.Players)+1 <= gameMaxPlayersCount {
		g.Players[p.Id] = p
		return nil
	}

	return errors.New("Game is full")
}

func (g *Game) Run() {
	ticker := time.NewTicker(time.Second / ticksPerSecond)
	defer ticker.Stop()

	for range ticker.C {
		if len(g.Players) == 2 {
			g.CheckPaddleCollision()
			g.Update()
			g.BroadcastGameState()
		}
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

func (g *Game) WaitForPlayers() {
	log.Printf("Players connected: %d", len(g.Players))
	for {
		if len(g.Players) == 2 {
			break
		}
	}
}

func (g *Game) HandlePaddleHit(racketX float64) {
	g.Ball.DY *= -1
}

func (g *Game) CheckPaddleCollision() {
	if g.Ball.Y >= g.BottomRacketY-racketHeight &&
		g.Ball.Y <= g.BottomRacketY &&
		g.Ball.X >= g.Players["player1"].X-(racketWidth/2) &&
		g.Ball.X <= g.Players["player1"].X+(racketWidth/2) {
		g.HandlePaddleHit(g.Players["player1"].X)
	}
	if g.Ball.Y <= g.TopRacketY+racketHeight &&
		g.Ball.Y >= g.TopRacketY &&
		g.Ball.X >= g.Players["player2"].X-(racketWidth/2) &&
		g.Ball.X <= g.Players["player2"].X+(racketWidth/2) {
		g.HandlePaddleHit(g.Players["player2"].X)
	}
}
