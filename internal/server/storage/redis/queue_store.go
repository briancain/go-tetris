package redis

import (
	"context"
	"time"
)

const queueKey = "matchmaking:queue"

// QueueStore implements Redis-based matchmaking queue
type QueueStore struct {
	client *Client
}

// NewQueueStore creates a new Redis queue store
func NewQueueStore(client *Client) *QueueStore {
	return &QueueStore{client: client}
}

// AddToQueue adds a player to the matchmaking queue
func (s *QueueStore) AddToQueue(playerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Remove player if already in queue, then add to end
	s.client.LRem(ctx, queueKey, 0, playerID)
	return s.client.RPush(ctx, queueKey, playerID).Err()
}

// RemoveFromQueue removes a player from the queue
func (s *QueueStore) RemoveFromQueue(playerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.client.LRem(ctx, queueKey, 0, playerID).Err()
}

// GetQueuedPlayers returns all players in the queue
func (s *QueueStore) GetQueuedPlayers() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.client.LRange(ctx, queueKey, 0, -1).Result()
}

// GetQueuePosition returns the position of a player in the queue (0-based)
func (s *QueueStore) GetQueuePosition(playerID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	players, err := s.client.LRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return -1, err
	}

	for i, id := range players {
		if id == playerID {
			return i, nil
		}
	}

	return -1, nil // Not in queue
}

// HealthCheck implements storage.HealthChecker
func (s *QueueStore) HealthCheck() error {
	return s.client.HealthCheck()
}
