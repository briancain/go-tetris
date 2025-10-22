package services

import (
	"testing"
	"time"

	"github.com/briancain/go-tetris/internal/server/storage/memory"
	"github.com/briancain/go-tetris/pkg/models"
)

func TestEndGame(t *testing.T) {
	// Setup
	gameStore := memory.NewGameStore()
	playerStore := memory.NewPlayerStore()
	wsManager := NewWebSocketManager()
	gm := NewGameManager(gameStore, playerStore, wsManager)

	// Create test players
	player1 := &models.Player{
		ID:       "gm_player1",
		Username: "testuser1",
		GameID:   "gm_game1",
	}
	player2 := &models.Player{
		ID:       "gm_player2",
		Username: "testuser2",
		GameID:   "gm_game1",
	}

	// Create test game
	game := &models.GameSession{
		ID:           "gm_game1",
		Player1:      player1,
		Player2:      player2,
		Player1Score: 1000,
		Player2Score: 800,
		Status:       models.GameStatusActive,
		CreatedAt:    time.Now(),
	}

	// Store game and players
	gameStore.CreateGame(game)
	playerStore.CreatePlayer(player1)
	playerStore.CreatePlayer(player2)

	// Test: Player 1 loses
	err := gm.EndGame("gm_game1", "gm_player1")
	if err != nil {
		t.Fatalf("EndGame failed: %v", err)
	}

	// Verify player 1 is marked as lost
	updatedGame, _ := gameStore.GetGame("gm_game1")
	if !updatedGame.Player1Lost {
		t.Error("Player1Lost should be true")
	}
	if updatedGame.Player2Lost {
		t.Error("Player2Lost should be false")
	}

	// Game should still be active (not finished yet)
	if updatedGame.Status == models.GameStatusFinished {
		t.Error("Game should not be finished when only one player loses")
	}
}

func TestCheckScoreWin(t *testing.T) {
	// Setup
	gameStore := memory.NewGameStore()
	playerStore := memory.NewPlayerStore()
	wsManager := NewWebSocketManager()
	gm := NewGameManager(gameStore, playerStore, wsManager)

	// Create test game where player1 lost with score 1000
	game := &models.GameSession{
		ID:           "gm_game2",
		Player1:      &models.Player{ID: "gm_player3", GameID: "gm_game2"},
		Player2:      &models.Player{ID: "gm_player4", GameID: "gm_game2"},
		Player1Lost:  true,
		Player2Lost:  false,
		Player1Score: 1000,
		Player2Score: 1200, // Player2 has higher score
		Status:       models.GameStatusActive,
	}

	gameStore.CreateGame(game)
	playerStore.CreatePlayer(game.Player1)
	playerStore.CreatePlayer(game.Player2)

	// Test score win detection
	gm.checkScoreWin(game)

	// Game should be finished since player2 beat player1's score
	updatedGame, _ := gameStore.GetGame("gm_game2")
	if updatedGame.Status != models.GameStatusFinished {
		t.Error("Game should be finished when surviving player beats loser's score")
	}
}

func TestCheckScoreWinNoWin(t *testing.T) {
	// Setup
	gameStore := memory.NewGameStore()
	playerStore := memory.NewPlayerStore()
	wsManager := NewWebSocketManager()
	gm := NewGameManager(gameStore, playerStore, wsManager)

	// Create test game where player1 lost but player2 hasn't beaten the score yet
	game := &models.GameSession{
		ID:           "gm_game3",
		Player1:      &models.Player{ID: "gm_player5", GameID: "gm_game3"},
		Player2:      &models.Player{ID: "gm_player6", GameID: "gm_game3"},
		Player1Lost:  true,
		Player2Lost:  false,
		Player1Score: 1000,
		Player2Score: 800, // Player2 has lower score
		Status:       models.GameStatusActive,
	}

	gameStore.CreateGame(game)

	// Test score win detection
	gm.checkScoreWin(game)

	// Game should still be active
	if game.Status != models.GameStatusActive {
		t.Error("Game should remain active when surviving player hasn't beaten loser's score")
	}
}

func TestBothPlayersLose(t *testing.T) {
	// Setup
	gameStore := memory.NewGameStore()
	playerStore := memory.NewPlayerStore()
	wsManager := NewWebSocketManager()
	gm := NewGameManager(gameStore, playerStore, wsManager)

	// Create test game
	game := &models.GameSession{
		ID:           "gm_game4",
		Player1:      &models.Player{ID: "gm_player7", GameID: "gm_game4"},
		Player2:      &models.Player{ID: "gm_player8", GameID: "gm_game4"},
		Player1Lost:  true, // Player1 already lost
		Player2Lost:  false,
		Player1Score: 1000,
		Player2Score: 800,
		Status:       models.GameStatusActive,
	}

	gameStore.CreateGame(game)
	playerStore.CreatePlayer(game.Player1)
	playerStore.CreatePlayer(game.Player2)

	// Player 2 also loses
	err := gm.EndGame("gm_game4", "gm_player8")
	if err != nil {
		t.Fatalf("EndGame failed: %v", err)
	}

	// Game should be finished since both players lost
	updatedGame, _ := gameStore.GetGame("gm_game4")
	if updatedGame.Status != models.GameStatusFinished {
		t.Error("Game should be finished when both players lose")
	}
	if !updatedGame.Player2Lost {
		t.Error("Player2Lost should be true")
	}
}

func TestHandleGameStateScoreTracking(t *testing.T) {
	// Setup
	gameStore := memory.NewGameStore()
	playerStore := memory.NewPlayerStore()
	wsManager := NewWebSocketManager()
	gm := NewGameManager(gameStore, playerStore, wsManager)

	// Create test players and game
	player1 := &models.Player{ID: "gm_player9", GameID: "gm_game5"}
	player2 := &models.Player{ID: "gm_player10", GameID: "gm_game5"}
	game := &models.GameSession{
		ID:      "gm_game5",
		Player1: player1,
		Player2: player2,
		Status:  models.GameStatusActive,
	}

	gameStore.CreateGame(game)
	playerStore.CreatePlayer(player1)
	playerStore.CreatePlayer(player2)

	// Test score update for player1
	gameState := &models.GameState{
		PlayerID: "gm_player9",
		GameID:   "gm_game5",
		Score:    1500,
		Level:    5,
		Lines:    20,
	}

	err := gm.HandleGameState("gm_player9", gameState)
	if err != nil {
		t.Fatalf("HandleGameState failed: %v", err)
	}

	// Verify score was updated
	updatedGame, _ := gameStore.GetGame("gm_game5")
	if updatedGame.Player1Score != 1500 {
		t.Errorf("Expected Player1Score to be 1500, got %d", updatedGame.Player1Score)
	}

	// Test score update for player2
	gameState.PlayerID = "gm_player10"
	gameState.Score = 1200

	err = gm.HandleGameState("gm_player10", gameState)
	if err != nil {
		t.Fatalf("HandleGameState failed: %v", err)
	}

	// Verify score was updated
	updatedGame, _ = gameStore.GetGame("gm_game5")
	if updatedGame.Player2Score != 1200 {
		t.Errorf("Expected Player2Score to be 1200, got %d", updatedGame.Player2Score)
	}
}
