package redis

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/briancain/go-tetris/pkg/models"
)

func TestGameStore_CreateGame(t *testing.T) {
	client := &Client{Client: redis.NewClient(&redis.Options{Addr: "localhost:6379"})}
	store := NewGameStore(client)

	// Create test game
	game := &models.GameSession{
		ID:        "test-game-1",
		Status:    models.GameStatusActive,
		Seed:      12345,
		CreatedAt: time.Now(),
	}

	err := store.CreateGame(game)
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
	}

	// Verify game was created
	retrieved, err := store.GetGame(game.ID)
	if err != nil {
		t.Fatalf("GetGame failed: %v", err)
	}

	if retrieved.ID != game.ID {
		t.Errorf("Expected ID %s, got %s", game.ID, retrieved.ID)
	}

	// Clean up
	store.DeleteGame(game.ID)
}

func TestGameStore_UpdateGame(t *testing.T) {
	client := &Client{Client: redis.NewClient(&redis.Options{Addr: "localhost:6379"})}
	store := NewGameStore(client)

	// Create test game
	game := &models.GameSession{
		ID:        "test-game-2",
		Status:    models.GameStatusActive,
		Seed:      12345,
		CreatedAt: time.Now(),
	}

	err := store.CreateGame(game)
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
	}

	// Update game status
	game.Status = models.GameStatusFinished
	err = store.UpdateGame(game)
	if err != nil {
		t.Fatalf("UpdateGame failed: %v", err)
	}

	// Verify update
	retrieved, err := store.GetGame(game.ID)
	if err != nil {
		t.Fatalf("GetGame failed: %v", err)
	}

	if retrieved.Status != models.GameStatusFinished {
		t.Errorf("Expected status %s, got %s", models.GameStatusFinished, retrieved.Status)
	}

	// Clean up
	store.DeleteGame(game.ID)
}

func TestGameStore_GetActiveGames(t *testing.T) {
	client := &Client{Client: redis.NewClient(&redis.Options{Addr: "localhost:6379"})}
	store := NewGameStore(client)

	// Create active game
	game := &models.GameSession{
		ID:        "test-game-3",
		Status:    models.GameStatusActive,
		Seed:      12345,
		CreatedAt: time.Now(),
	}

	err := store.CreateGame(game)
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
	}

	// Get active games
	activeGames, err := store.GetActiveGames()
	if err != nil {
		t.Fatalf("GetActiveGames failed: %v", err)
	}

	// Should contain our game
	found := false
	for _, g := range activeGames {
		if g.ID == game.ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("Active game not found in GetActiveGames result")
	}

	// Clean up
	store.DeleteGame(game.ID)
}

func TestGameStore_HealthCheck(t *testing.T) {
	client := &Client{Client: redis.NewClient(&redis.Options{Addr: "localhost:6379"})}
	store := NewGameStore(client)

	err := store.HealthCheck()
	if err != nil {
		t.Skipf("Redis not available for testing: %v", err)
	}
}
