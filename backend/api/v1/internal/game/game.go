package game

import (
	"crypto/rand"
	"errors"
	"sync"
)

func InitGame() *Game {
	return &Game{
		Id:      generateGameId(),
		Players: make(map[string]*Player),
		Ball: Ball{
			X:  0.5,
			Y:  0.5,
			DX: 0,
			DY: 0,
		},
		Started: false,
		mu:      &sync.Mutex{},
	}
}

func (g *Game) AddPlayer(p *Player) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i := range g.Players {
		if g.Players[i] == nil {
			g.Players[i] = p
			return nil
		}
	}

	return errors.New("game is full")
}

func generateGameId() string {
	return rand.Text()
}
