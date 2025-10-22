package memory

import (
	"testing"

	"github.com/briancain/go-tetris/pkg/models"
)

func TestPlayerStore_CreateAndGet(t *testing.T) {
	store := NewPlayerStore()

	player := &models.Player{
		ID:           "test-id",
		Username:     "testuser",
		SessionToken: "test-token",
	}

	// Test create
	err := store.CreatePlayer(player)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test get by ID
	retrieved, err := store.GetPlayer("test-id")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", retrieved.Username)
	}

	// Test get by token
	retrievedByToken, err := store.GetPlayerByToken("test-token")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrievedByToken.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", retrievedByToken.ID)
	}
}

func TestPlayerStore_CreateDuplicate(t *testing.T) {
	store := NewPlayerStore()

	player := &models.Player{
		ID:           "test-id",
		Username:     "testuser",
		SessionToken: "test-token",
	}

	// Create first time
	err := store.CreatePlayer(player)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Try to create duplicate
	err = store.CreatePlayer(player)
	if err == nil {
		t.Error("Expected error for duplicate player")
	}
}

func TestPlayerStore_UpdatePlayer(t *testing.T) {
	store := NewPlayerStore()

	player := &models.Player{
		ID:           "test-id",
		Username:     "testuser",
		SessionToken: "test-token",
		InQueue:      false,
	}

	// Create player
	store.CreatePlayer(player)

	// Update player
	player.InQueue = true
	err := store.UpdatePlayer(player)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update
	retrieved, _ := store.GetPlayer("test-id")
	if !retrieved.InQueue {
		t.Error("Expected InQueue to be true after update")
	}
}

func TestPlayerStore_DeletePlayer(t *testing.T) {
	store := NewPlayerStore()

	player := &models.Player{
		ID:           "test-id",
		Username:     "testuser",
		SessionToken: "test-token",
	}

	// Create player
	store.CreatePlayer(player)

	// Delete player
	err := store.DeletePlayer("test-id")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = store.GetPlayer("test-id")
	if err == nil {
		t.Error("Expected error for deleted player")
	}

	// Verify token is also removed
	_, err = store.GetPlayerByToken("test-token")
	if err == nil {
		t.Error("Expected error for deleted player token")
	}
}
