package memory

import (
	"sync"
)

// QueueStore implements in-memory matchmaking queue
type QueueStore struct {
	queue []string
	mu    sync.RWMutex
}

// NewQueueStore creates a new in-memory queue store
func NewQueueStore() *QueueStore {
	return &QueueStore{
		queue: make([]string, 0),
	}
}

// AddToQueue adds a player to the matchmaking queue
func (s *QueueStore) AddToQueue(playerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if player is already in queue
	for _, id := range s.queue {
		if id == playerID {
			return nil // Already in queue
		}
	}

	s.queue = append(s.queue, playerID)
	return nil
}

// RemoveFromQueue removes a player from the queue
func (s *QueueStore) RemoveFromQueue(playerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, id := range s.queue {
		if id == playerID {
			s.queue = append(s.queue[:i], s.queue[i+1:]...)
			return nil
		}
	}

	return nil // Not in queue, no error
}

// GetQueuedPlayers returns all players in the queue
func (s *QueueStore) GetQueuedPlayers() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]string, len(s.queue))
	copy(result, s.queue)
	return result, nil
}

// GetQueuePosition returns the position of a player in the queue (0-based)
func (s *QueueStore) GetQueuePosition(playerID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, id := range s.queue {
		if id == playerID {
			return i, nil
		}
	}

	return -1, nil // Not in queue
}
