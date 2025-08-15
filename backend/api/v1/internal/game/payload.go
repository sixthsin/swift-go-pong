package game

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Conn  *websocket.Conn
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	X     float64 `json:"x"`
	Score int     `json:"score"`
}

type Ball struct {
	Speed float64 `json:"speed"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	DX    float64 `json:"dx"`
	DY    float64 `json:"dy"`
}

type Game struct {
	Id      string     `json:"id"`
	Players [2]*Player `json:"players"`
	Ball    Ball       `json:"ball"`
	started bool
	mu      sync.Mutex
}
