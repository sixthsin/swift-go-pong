package game

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Conn  *websocket.Conn
	Id    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Score int     `json:"score"`
}

type Ball struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	DX float64 `json:"dx"`
	DY float64 `json:"dy"`
}

type Game struct {
	Id            string             `json:"id"`
	Players       map[string]*Player `json:"players"`
	Ball          Ball               `json:"ball"`
	Started       bool               `json:"started"`
	Mu            *sync.Mutex        `json:"-"`
	TopRacketY    float64            `json:"-"`
	BottomRacketY float64            `json:"-"`
}

type PlayerInput struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
}
