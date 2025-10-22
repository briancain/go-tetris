package memory

import (
	"testing"
)

func TestQueueStore_AddAndGet(t *testing.T) {
	store := NewQueueStore()

	// Add players to queue
	err := store.AddToQueue("player1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = store.AddToQueue("player2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Get queued players
	players, err := store.GetQueuedPlayers()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(players))
	}

	if players[0] != "player1" || players[1] != "player2" {
		t.Errorf("Expected [player1, player2], got %v", players)
	}
}

func TestQueueStore_GetPosition(t *testing.T) {
	store := NewQueueStore()

	// Add players
	store.AddToQueue("player1")
	store.AddToQueue("player2")
	store.AddToQueue("player3")

	// Check positions
	pos1, err := store.GetQueuePosition("player1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if pos1 != 0 {
		t.Errorf("Expected position 0 for player1, got %d", pos1)
	}

	pos2, err := store.GetQueuePosition("player2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if pos2 != 1 {
		t.Errorf("Expected position 1 for player2, got %d", pos2)
	}

	// Check non-existent player
	pos, err := store.GetQueuePosition("nonexistent")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if pos != -1 {
		t.Errorf("Expected position -1 for non-existent player, got %d", pos)
	}
}

func TestQueueStore_RemoveFromQueue(t *testing.T) {
	store := NewQueueStore()

	// Add players
	store.AddToQueue("player1")
	store.AddToQueue("player2")
	store.AddToQueue("player3")

	// Remove middle player
	err := store.RemoveFromQueue("player2")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check remaining players
	players, _ := store.GetQueuedPlayers()
	if len(players) != 2 {
		t.Errorf("Expected 2 players after removal, got %d", len(players))
	}

	if players[0] != "player1" || players[1] != "player3" {
		t.Errorf("Expected [player1, player3], got %v", players)
	}

	// Check positions updated
	pos3, _ := store.GetQueuePosition("player3")
	if pos3 != 1 {
		t.Errorf("Expected player3 to be at position 1, got %d", pos3)
	}
}

func TestQueueStore_AddDuplicate(t *testing.T) {
	store := NewQueueStore()

	// Add player twice
	store.AddToQueue("player1")
	store.AddToQueue("player1")

	// Should only be in queue once
	players, _ := store.GetQueuedPlayers()
	if len(players) != 1 {
		t.Errorf("Expected 1 player (no duplicates), got %d", len(players))
	}
}
