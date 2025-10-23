package redis

import (
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestQueueStore_AddToQueue(t *testing.T) {
	// Use Redis mock for unit testing
	client := &Client{Client: redis.NewClient(&redis.Options{Addr: "localhost:6379"})}
	store := NewQueueStore(client)

	// Test adding player to queue
	err := store.AddToQueue("player1")
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
	}

	// Test getting queue position
	pos, err := store.GetQueuePosition("player1")
	if err != nil {
		t.Fatalf("GetQueuePosition failed: %v", err)
	}
	if pos != 0 {
		t.Errorf("Expected position 0, got %d", pos)
	}

	// Clean up
	store.RemoveFromQueue("player1")
}

func TestQueueStore_GetQueuedPlayers(t *testing.T) {
	client := &Client{Client: redis.NewClient(&redis.Options{Addr: "localhost:6379"})}
	store := NewQueueStore(client)

	// Test with empty queue
	players, err := store.GetQueuedPlayers()
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
	}

	// Add players
	store.AddToQueue("player1")
	store.AddToQueue("player2")

	players, err = store.GetQueuedPlayers()
	if err != nil {
		t.Fatalf("GetQueuedPlayers failed: %v", err)
	}

	if len(players) < 2 {
		t.Errorf("Expected at least 2 players, got %d", len(players))
	}

	// Clean up
	store.RemoveFromQueue("player1")
	store.RemoveFromQueue("player2")
}

func TestQueueStore_HealthCheck(t *testing.T) {
	client := &Client{Client: redis.NewClient(&redis.Options{Addr: "localhost:6379"})}
	store := NewQueueStore(client)

	err := store.HealthCheck()
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
	}
}
