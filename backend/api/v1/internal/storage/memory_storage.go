package storage

import (
	"errors"
	"pong-api-v1/internal/game"
	"sync"
)

type MemoryStorage struct {
	games map[string]*game.Game
	mu    sync.Mutex
}

func NewStorage() *MemoryStorage {
	return &MemoryStorage{
		games: make(map[string]*game.Game),
		mu:    sync.Mutex{},
	}
}

func (s *MemoryStorage) Save(newGame *game.Game) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.games[newGame.Id] = newGame
}

func (s *MemoryStorage) Get(gameId string) (*game.Game, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	game, exists := s.games[gameId]
	if !exists {
		return nil, errors.New("game not found")
	}

	return game, nil
}

func (s *MemoryStorage) Delete(gameId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.games[gameId]
	if !exists {
		return errors.New("game not found")
	}

	delete(s.games, gameId)

	return nil
}
