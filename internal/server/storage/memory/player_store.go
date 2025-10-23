package memory

import (
	"errors"
	"sync"

	"github.com/briancain/go-tetris/pkg/models"
)

// PlayerStore implements in-memory player storage
type PlayerStore struct {
	players map[string]*models.Player
	tokens  map[string]string // token -> playerID
	mu      sync.RWMutex
}

// NewPlayerStore creates a new in-memory player store
func NewPlayerStore() *PlayerStore {
	return &PlayerStore{
		players: make(map[string]*models.Player),
		tokens:  make(map[string]string),
	}
}

// CreatePlayer stores a new player
func (s *PlayerStore) CreatePlayer(player *models.Player) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.players[player.ID]; exists {
		return errors.New("player already exists")
	}

	s.players[player.ID] = player
	s.tokens[player.SessionToken] = player.ID
	return nil
}

// GetPlayer retrieves a player by ID
func (s *PlayerStore) GetPlayer(id string) (*models.Player, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	player, exists := s.players[id]
	if !exists {
		return nil, errors.New("player not found")
	}

	return player, nil
}

// GetPlayerByUsername retrieves a player by username
func (s *PlayerStore) GetPlayerByUsername(username string) (*models.Player, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, player := range s.players {
		if player.Username == username {
			return player, nil
		}
	}

	return nil, errors.New("player not found")
}

// GetPlayerByToken retrieves a player by session token
func (s *PlayerStore) GetPlayerByToken(token string) (*models.Player, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	playerID, exists := s.tokens[token]
	if !exists {
		return nil, errors.New("invalid token")
	}

	return s.GetPlayer(playerID)
}

// UpdatePlayer updates an existing player
func (s *PlayerStore) UpdatePlayer(player *models.Player) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.players[player.ID]; !exists {
		return errors.New("player not found")
	}

	s.players[player.ID] = player
	return nil
}

// DeletePlayer removes a player
func (s *PlayerStore) DeletePlayer(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, exists := s.players[id]
	if !exists {
		return errors.New("player not found")
	}

	delete(s.players, id)
	delete(s.tokens, player.SessionToken)
	return nil
}

// GetAllPlayers returns all players
func (s *PlayerStore) GetAllPlayers() ([]*models.Player, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var players []*models.Player
	for _, player := range s.players {
		players = append(players, player)
	}

	return players, nil
}
