package services

import (
	"sync"
	"testing"

	"github.com/briancain/go-tetris/internal/server/storage/memory"
	"github.com/briancain/go-tetris/pkg/models"
)

// TestGameManagerConcurrentAccess tests that GameManager handles concurrent operations safely
func TestGameManagerConcurrentAccess(t *testing.T) {
	// Setup
	gameStore := memory.NewGameStore()
	playerStore := memory.NewPlayerStore()
	wsManager := NewWebSocketManager()
	gm := NewGameManager(gameStore, playerStore, wsManager)

	// Create test players
	player1 := &models.Player{
		ID:       "player1",
		Username: "testplayer1",
	}
	player2 := &models.Player{
		ID:       "player2",
		Username: "testplayer2",
	}

	err := playerStore.CreatePlayer(player1)
	if err != nil {
		t.Fatalf("Failed to create player1: %v", err)
	}
	err = playerStore.CreatePlayer(player2)
	if err != nil {
		t.Fatalf("Failed to create player2: %v", err)
	}

	// Create test game
	game := &models.GameSession{
		ID:      "testgame",
		Player1: player1,
		Player2: player2,
		Status:  models.GameStatusActive,
		Seed:    12345,
	}
	player1.GameID = game.ID
	player2.GameID = game.ID

	err = gameStore.CreateGame(game)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	err = playerStore.UpdatePlayer(player1)
	if err != nil {
		t.Fatalf("Failed to update player1: %v", err)
	}
	err = playerStore.UpdatePlayer(player2)
	if err != nil {
		t.Fatalf("Failed to update player2: %v", err)
	}

	// Test concurrent operations
	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // 3 types of operations

	// Concurrent game moves
	for i := 0; i < numGoroutines; i++ {
		go func(playerID string) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				move := &models.GameMove{
					PlayerID: playerID,
					GameID:   game.ID,
					MoveType: "left",
				}
				gm.HandleGameMove(playerID, move)
			}
		}(player1.ID)
	}

	// Concurrent game state updates
	for i := 0; i < numGoroutines; i++ {
		go func(playerID string) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				state := &models.GameState{
					Score: j,
					Level: 1,
					Lines: j,
				}
				gm.HandleGameState(playerID, state)
			}
		}(player2.ID)
	}

	// Concurrent disconnect attempts
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				// Alternate between players
				if j%2 == 0 {
					gm.HandlePlayerDisconnect(player1.ID)
				} else {
					gm.HandlePlayerDisconnect(player2.ID)
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// If we reach here without panicking, the race condition is fixed
	t.Log("GameManager handled concurrent operations without race conditions")
}

// TestGameManagerConcurrentEndGame tests concurrent game ending scenarios
func TestGameManagerConcurrentEndGame(t *testing.T) {
	// Setup
	gameStore := memory.NewGameStore()
	playerStore := memory.NewPlayerStore()
	wsManager := NewWebSocketManager()
	gm := NewGameManager(gameStore, playerStore, wsManager)

	// Create multiple games for concurrent ending
	const numGames = 50
	var wg sync.WaitGroup
	wg.Add(numGames)

	for i := 0; i < numGames; i++ {
		go func(gameIndex int) {
			defer wg.Done()

			// Create unique players and game for each goroutine
			player1 := &models.Player{
				ID:       "player1_" + string(rune('0'+gameIndex)),
				Username: "testplayer1_" + string(rune('0'+gameIndex)),
			}
			player2 := &models.Player{
				ID:       "player2_" + string(rune('0'+gameIndex)),
				Username: "testplayer2_" + string(rune('0'+gameIndex)),
			}

			playerStore.CreatePlayer(player1)
			playerStore.CreatePlayer(player2)

			game := &models.GameSession{
				ID:      "testgame_" + string(rune('0'+gameIndex)),
				Player1: player1,
				Player2: player2,
				Status:  models.GameStatusActive,
				Seed:    int64(12345 + gameIndex),
			}
			player1.GameID = game.ID
			player2.GameID = game.ID

			gameStore.CreateGame(game)
			playerStore.UpdatePlayer(player1)
			playerStore.UpdatePlayer(player2)

			// Try to end the game (this should be safe with mutex protection)
			gm.EndGame(game.ID, player1.ID)
		}(i)
	}

	// Wait for all games to end
	wg.Wait()

	t.Log("GameManager handled concurrent game endings without race conditions")
}
