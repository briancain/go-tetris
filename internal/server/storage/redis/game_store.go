package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/briancain/go-tetris/pkg/models"
)

const (
	gameKeyPrefix  = "game:"
	activeGamesKey = "games:active"
	allGamesKey    = "games:all"
	gameSessionTTL = 2 * time.Hour // Games expire after 2 hours of inactivity
)

// GameStore implements Redis-based game session storage
type GameStore struct {
	client *Client
}

// NewGameStore creates a new Redis game store
func NewGameStore(client *Client) *GameStore {
	return &GameStore{client: client}
}

// CreateGame stores a new game session
func (s *GameStore) CreateGame(game *models.GameSession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameKey := gameKeyPrefix + game.ID

	// Check if game already exists
	exists, err := s.client.Exists(ctx, gameKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("game already exists")
	}

	// Serialize game to JSON
	data, err := json.Marshal(game)
	if err != nil {
		return err
	}

	// Store game with TTL
	err = s.client.Set(ctx, gameKey, data, gameSessionTTL).Err()
	if err != nil {
		return err
	}

	// Add to active games set
	err = s.client.SAdd(ctx, activeGamesKey, game.ID).Err()
	if err != nil {
		return err
	}

	// Add to all games set
	return s.client.SAdd(ctx, allGamesKey, game.ID).Err()
}

// GetGame retrieves a game session by ID
func (s *GameStore) GetGame(id string) (*models.GameSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameKey := gameKeyPrefix + id

	data, err := s.client.Get(ctx, gameKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, errors.New("game not found")
		}
		return nil, err
	}

	var game models.GameSession
	err = json.Unmarshal([]byte(data), &game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

// UpdateGame updates an existing game session
func (s *GameStore) UpdateGame(game *models.GameSession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameKey := gameKeyPrefix + game.ID

	// Check if game exists
	exists, err := s.client.Exists(ctx, gameKey).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("game not found")
	}

	// Serialize game to JSON
	data, err := json.Marshal(game)
	if err != nil {
		return err
	}

	// Update game with TTL refresh
	err = s.client.Set(ctx, gameKey, data, gameSessionTTL).Err()
	if err != nil {
		return err
	}

	// Update active games set based on status
	if game.Status == models.GameStatusFinished {
		return s.client.SRem(ctx, activeGamesKey, game.ID).Err()
	}

	return nil
}

// DeleteGame removes a game session
func (s *GameStore) DeleteGame(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameKey := gameKeyPrefix + id

	// Delete game data
	err := s.client.Del(ctx, gameKey).Err()
	if err != nil {
		return err
	}

	// Remove from sets
	s.client.SRem(ctx, activeGamesKey, id)
	s.client.SRem(ctx, allGamesKey, id)

	return nil
}

// GetActiveGames returns all active game sessions
func (s *GameStore) GetActiveGames() ([]*models.GameSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get active game IDs
	gameIDs, err := s.client.SMembers(ctx, activeGamesKey).Result()
	if err != nil {
		return nil, err
	}

	var games []*models.GameSession
	for _, id := range gameIDs {
		game, err := s.GetGame(id)
		if err != nil {
			// Game might have expired, remove from active set
			s.client.SRem(ctx, activeGamesKey, id)
			continue
		}
		games = append(games, game)
	}

	return games, nil
}

// GetAllGames returns all game sessions
func (s *GameStore) GetAllGames() ([]*models.GameSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get all game IDs
	gameIDs, err := s.client.SMembers(ctx, allGamesKey).Result()
	if err != nil {
		return nil, err
	}

	var games []*models.GameSession
	for _, id := range gameIDs {
		game, err := s.GetGame(id)
		if err != nil {
			// Game might have expired, remove from all games set
			s.client.SRem(ctx, allGamesKey, id)
			continue
		}
		games = append(games, game)
	}

	return games, nil
}

// HealthCheck implements storage.HealthChecker
func (s *GameStore) HealthCheck() error {
	return s.client.HealthCheck()
}
