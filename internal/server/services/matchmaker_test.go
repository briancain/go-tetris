package services

import (
	"testing"
	"time"

	"github.com/briancain/go-tetris/internal/server/storage/memory"
	"github.com/briancain/go-tetris/pkg/models"
)

func TestMatchmakingService_JoinQueue(t *testing.T) {
	// Setup
	playerStore := memory.NewPlayerStore()
	gameStore := memory.NewGameStore()
	queueStore := memory.NewQueueStore()
	wsManager := NewWebSocketManager()
	gameManager := NewGameManager(gameStore, playerStore, wsManager)
	matchmaker := NewMatchmakingService(playerStore, gameStore, queueStore, gameManager)

	// Create a player
	player := &models.Player{
		ID:       "player1",
		Username: "testuser",
		InQueue:  false,
	}
	playerStore.CreatePlayer(player)

	// Test joining queue
	err := matchmaker.JoinQueue("player1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify player is in queue
	position, err := matchmaker.GetQueueStatus("player1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if position != 0 {
		t.Errorf("Expected position 0, got %d", position)
	}
}

func TestMatchmakingService_LeaveQueue(t *testing.T) {
	// Setup
	playerStore := memory.NewPlayerStore()
	gameStore := memory.NewGameStore()
	queueStore := memory.NewQueueStore()
	wsManager := NewWebSocketManager()
	gameManager := NewGameManager(gameStore, playerStore, wsManager)
	matchmaker := NewMatchmakingService(playerStore, gameStore, queueStore, gameManager)

	// Create a player
	player := &models.Player{
		ID:       "player1",
		Username: "testuser",
		InQueue:  false,
	}
	playerStore.CreatePlayer(player)

	// Join then leave queue
	matchmaker.JoinQueue("player1")
	err := matchmaker.LeaveQueue("player1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify player is not in queue
	position, err := matchmaker.GetQueueStatus("player1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if position != -1 {
		t.Errorf("Expected position -1 (not in queue), got %d", position)
	}
}

func TestMatchmakingService_TwoPlayerMatch(t *testing.T) {
	// Setup with fresh stores
	playerStore := memory.NewPlayerStore()
	gameStore := memory.NewGameStore()
	queueStore := memory.NewQueueStore()
	wsManager := NewWebSocketManager()
	gameManager := NewGameManager(gameStore, playerStore, wsManager)
	matchmaker := NewMatchmakingService(playerStore, gameStore, queueStore, gameManager)

	// Create two players
	player1 := &models.Player{
		ID:       "player1",
		Username: "user1",
		InQueue:  false,
	}
	player2 := &models.Player{
		ID:       "player2",
		Username: "user2",
		InQueue:  false,
	}
	playerStore.CreatePlayer(player1)
	playerStore.CreatePlayer(player2)

	// Both join queue
	matchmaker.JoinQueue("player1")
	matchmaker.JoinQueue("player2")

	// Wait for matchmaking to complete with longer timeout and more checks
	var games []*models.GameSession
	var err error
	matched := false

	for i := 0; i < 20; i++ { // Increased iterations
		time.Sleep(25 * time.Millisecond) // Shorter sleep, more frequent checks

		// Check if players are matched (out of queue)
		pos1, _ := matchmaker.GetQueueStatus("player1")
		pos2, _ := matchmaker.GetQueueStatus("player2")

		if pos1 == -1 && pos2 == -1 {
			matched = true
			break
		}
	}

	if !matched {
		t.Fatal("Players were not matched within timeout")
	}

	// Now check for exactly one game
	games, err = gameStore.GetActiveGames()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(games) != 1 {
		t.Errorf("Expected exactly 1 active game, got %d", len(games))
		return
	}

	game := games[0]
	if game.Player1.ID != "player1" || game.Player2.ID != "player2" {
		t.Errorf("Expected players player1 and player2, got %s and %s",
			game.Player1.ID, game.Player2.ID)
	}

	if game.Seed == 0 {
		t.Error("Expected game seed to be generated")
	}
}
