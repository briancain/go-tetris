package memory

import (
	"errors"
	"sync"

	"github.com/briancain/go-tetris/pkg/models"
)

// GameStore implements in-memory game session storage
type GameStore struct {
	games map[string]*models.GameSession
	mu    sync.RWMutex
}

// NewGameStore creates a new in-memory game store
func NewGameStore() *GameStore {
	return &GameStore{
		games: make(map[string]*models.GameSession),
	}
}

// CreateGame stores a new game session
func (s *GameStore) CreateGame(game *models.GameSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.games[game.ID]; exists {
		return errors.New("game already exists")
	}

	s.games[game.ID] = game
	return nil
}

// GetGame retrieves a game session by ID
func (s *GameStore) GetGame(id string) (*models.GameSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	game, exists := s.games[id]
	if !exists {
		return nil, errors.New("game not found")
	}

	return game, nil
}

// UpdateGame updates an existing game session
func (s *GameStore) UpdateGame(game *models.GameSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.games[game.ID]; !exists {
		return errors.New("game not found")
	}

	s.games[game.ID] = game
	return nil
}

// DeleteGame removes a game session
func (s *GameStore) DeleteGame(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.games[id]; !exists {
		return errors.New("game not found")
	}

	delete(s.games, id)
	return nil
}

// GetActiveGames returns all active game sessions
func (s *GameStore) GetActiveGames() ([]*models.GameSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var activeGames []*models.GameSession
	for _, game := range s.games {
		if game.Status == models.GameStatusActive || game.Status == models.GameStatusWaiting {
			activeGames = append(activeGames, game)
		}
	}

	return activeGames, nil
}

// GetAllGames returns all game sessions
func (s *GameStore) GetAllGames() ([]*models.GameSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allGames []*models.GameSession
	for _, game := range s.games {
		allGames = append(allGames, game)
	}

	return allGames, nil
}
