package tetris

import (
	"testing"

	_ "github.com/briancain/go-tetris/internal/testutil" // Import for init side effects
)

func TestNewGame(t *testing.T) {
	game := NewGame()

	if game == nil {
		t.Fatal("NewGame returned nil")
	}

	if game.Board == nil {
		t.Error("Board is nil")
	}

	if game.State != StateMenu {
		t.Errorf("Expected initial state to be StateMenu, got %d", game.State)
	}

	if game.Score != 0 {
		t.Errorf("Expected initial score to be 0, got %d", game.Score)
	}

	if game.Level != 1 {
		t.Errorf("Expected initial level to be 1, got %d", game.Level)
	}

	if game.LinesCleared != 0 {
		t.Errorf("Expected initial lines cleared to be 0, got %d", game.LinesCleared)
	}
}

func TestStart(t *testing.T) {
	game := NewGame()
	game.Start()

	if game.State != StatePlaying {
		t.Errorf("Expected game state to be StatePlaying, got %d", game.State)
	}

	if game.CurrentPiece == nil {
		t.Error("CurrentPiece is nil after starting")
	}

	if game.NextPiece == nil {
		t.Error("NextPiece is nil after starting")
	}
}

func TestGameState(t *testing.T) {
	game := NewGame()
	game.Start()

	// Test game state changes
	// Set to paused
	game.State = StatePaused
	if game.State != StatePaused {
		t.Errorf("Expected game state to be StatePaused, got %d", game.State)
	}

	// Set back to playing
	game.State = StatePlaying
	if game.State != StatePlaying {
		t.Errorf("Expected game state to be StatePlaying, got %d", game.State)
	}

	// Set to game over
	game.State = StateGameOver
	if game.State != StateGameOver {
		t.Errorf("Expected game state to be StateGameOver, got %d", game.State)
	}
}

func TestGetScore(t *testing.T) {
	game := NewGame()
	game.Score = 100

	if game.GetScore() != 100 {
		t.Errorf("Expected score to be 100, got %d", game.GetScore())
	}
}

func TestGetLevel(t *testing.T) {
	game := NewGame()
	game.Level = 5

	if game.GetLevel() != 5 {
		t.Errorf("Expected level to be 5, got %d", game.GetLevel())
	}
}

func TestGetLinesCleared(t *testing.T) {
	game := NewGame()
	game.LinesCleared = 10

	if game.GetLinesCleared() != 10 {
		t.Errorf("Expected lines cleared to be 10, got %d", game.GetLinesCleared())
	}
}

func TestHoldPiece(t *testing.T) {
	game := NewGame()
	game.Start()

	// Initially, held piece should be nil
	if game.HeldPiece != nil {
		t.Error("HeldPiece should be nil at game start")
	}

	// Hold the current piece
	currentPiece := game.CurrentPiece
	game.HoldPiece()

	// Now held piece should not be nil
	if game.HeldPiece == nil {
		t.Error("HeldPiece is nil after holding")
	}

	// Current piece should be different
	if game.CurrentPiece == currentPiece {
		t.Error("CurrentPiece should change after holding")
	}

	// HasSwapped should be true
	if !game.HasSwapped {
		t.Error("HasSwapped should be true after holding")
	}

	// Try to hold again, should not work
	heldPiece := game.HeldPiece
	game.HoldPiece()

	// Held piece should not change
	if game.HeldPiece != heldPiece {
		t.Error("HeldPiece should not change when holding twice in a row")
	}
}
