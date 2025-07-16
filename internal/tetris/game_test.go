package tetris

import (
	"testing"
	"time"
)

func TestNewGame(t *testing.T) {
	game := NewGame()

	// Check initial state
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

	if game.NextPiece == nil {
		t.Error("NextPiece is nil")
	}
}

func TestStart(t *testing.T) {
	game := NewGame()
	game.Start()

	// Check game state after starting
	if game.State != StatePlaying {
		t.Errorf("Expected game state to be StatePlaying, got %d", game.State)
	}

	if game.CurrentPiece == nil {
		t.Error("CurrentPiece is nil after starting")
	}

	if game.NextPiece == nil {
		t.Error("NextPiece is nil after starting")
	}

	if game.HeldPiece != nil {
		t.Error("HeldPiece should be nil at game start")
	}

	if game.HasSwapped {
		t.Error("HasSwapped should be false at game start")
	}
}

func TestMoveLeftRight(t *testing.T) {
	game := NewGame()
	game.Start()

	// Force a specific piece type and position for consistent testing
	game.CurrentPiece = NewPiece(TypeI)
	game.CurrentPiece.X = 3

	// Set LastMoveSide to a time far in the past to bypass input delay
	game.LastMoveSide = time.Time{}

	// Move left
	result := game.MoveLeft()
	if !result {
		t.Error("MoveLeft failed")
	}

	if game.CurrentPiece.X != 2 {
		t.Errorf("Expected X position 2, got %d", game.CurrentPiece.X)
	}

	// Reset position and LastMoveSide
	game.CurrentPiece.X = 3
	game.LastMoveSide = time.Time{}

	// Move right
	result = game.MoveRight()
	if !result {
		t.Error("MoveRight failed")
	}

	if game.CurrentPiece.X != 4 {
		t.Errorf("Expected X position 4, got %d", game.CurrentPiece.X)
	}
}

func TestRotatePiece(t *testing.T) {
	game := NewGame()
	game.Start()

	// Force a specific piece type for consistent testing
	game.CurrentPiece = NewPiece(TypeT)
	originalShape := make([][]bool, len(game.CurrentPiece.Shape))
	for i := range game.CurrentPiece.Shape {
		originalShape[i] = make([]bool, len(game.CurrentPiece.Shape[i]))
		copy(originalShape[i], game.CurrentPiece.Shape[i])
	}

	// Rotate the piece
	result := game.RotatePiece()
	if !result {
		t.Error("RotatePiece failed")
	}

	// Check that the shape changed
	shapeDifferent := false
	for i := range game.CurrentPiece.Shape {
		for j := range game.CurrentPiece.Shape[i] {
			if i < len(originalShape) && j < len(originalShape[i]) {
				if game.CurrentPiece.Shape[i][j] != originalShape[i][j] {
					shapeDifferent = true
				}
			} else {
				shapeDifferent = true
			}
		}
	}

	if !shapeDifferent {
		t.Error("Piece shape did not change after rotation")
	}
}

func TestHardDrop(t *testing.T) {
	game := NewGame()
	game.Start()

	// Set up a specific scenario
	game.CurrentPiece = NewPiece(TypeI)
	game.NextPiece = NewPiece(TypeO)

	// Hard drop
	game.HardDrop()

	// Check that the next piece became the current piece
	if game.CurrentPiece.Type != TypeO {
		t.Errorf("Expected current piece to be TypeO after hard drop, got %v", game.CurrentPiece.Type)
	}
}

func TestSoftDrop(t *testing.T) {
	game := NewGame()
	game.Start()

	// Save initial position
	initialY := game.CurrentPiece.Y

	// Soft drop
	result := game.SoftDrop()
	if !result {
		t.Error("SoftDrop failed")
	}

	if game.CurrentPiece.Y != initialY+1 {
		t.Errorf("Expected Y position %d, got %d", initialY+1, game.CurrentPiece.Y)
	}
}

func TestTogglePause(t *testing.T) {
	game := NewGame()
	game.Start()

	// Initially playing
	if game.State != StatePlaying {
		t.Errorf("Expected initial state to be StatePlaying, got %d", game.State)
	}

	// Pause
	game.TogglePause()
	if game.State != StatePaused {
		t.Errorf("Expected state to be StatePaused, got %d", game.State)
	}

	// Unpause
	game.TogglePause()
	if game.State != StatePlaying {
		t.Errorf("Expected state to be StatePlaying, got %d", game.State)
	}
}

func TestHoldPiece(t *testing.T) {
	game := NewGame()
	game.Start()

	// Initially no held piece
	if game.HeldPiece != nil {
		t.Error("Expected no held piece initially")
	}

	// Save current piece type
	currentType := game.CurrentPiece.Type
	nextType := game.NextPiece.Type

	// Hold the piece
	result := game.HoldPiece()
	if !result {
		t.Error("HoldPiece failed")
	}

	// Check that the held piece is the original current piece
	if game.HeldPiece == nil {
		t.Error("HeldPiece is nil after holding")
	} else if game.HeldPiece.Type != currentType {
		t.Errorf("Expected held piece type %v, got %v", currentType, game.HeldPiece.Type)
	}

	// Check that the current piece is now the next piece
	if game.CurrentPiece.Type != nextType {
		t.Errorf("Expected current piece to be the next piece type %v, got %v", nextType, game.CurrentPiece.Type)
	}

	// Check that HasSwapped is true
	if !game.HasSwapped {
		t.Error("HasSwapped should be true after holding")
	}

	// Try to hold again (should fail)
	result = game.HoldPiece()
	if result {
		t.Error("HoldPiece succeeded when it should have failed")
	}

	// Lock the piece and check that HasSwapped is reset
	game.lockPiece()
	if game.HasSwapped {
		t.Error("HasSwapped should be false after locking piece")
	}
}

func TestUpdateDropInterval(t *testing.T) {
	game := NewGame()
	initialInterval := game.DropInterval

	// Increase level
	game.Level = 2
	game.updateDropInterval()

	// Check that the interval decreased
	if game.DropInterval >= initialInterval {
		t.Errorf("Expected drop interval to decrease, got %v >= %v", game.DropInterval, initialInterval)
	}

	// Increase level further
	game.Level = 5
	previousInterval := game.DropInterval
	game.updateDropInterval()

	// Check that the interval decreased further
	if game.DropInterval >= previousInterval {
		t.Errorf("Expected drop interval to decrease further, got %v >= %v", game.DropInterval, previousInterval)
	}
}

func TestAddScore(t *testing.T) {
	game := NewGame()
	game.Level = 1

	testCases := []struct {
		linesCleared  int
		expectedScore int
	}{
		{1, 40},   // 1 line = 40 points at level 1
		{2, 100},  // 2 lines = 100 points at level 1
		{3, 300},  // 3 lines = 300 points at level 1
		{4, 1200}, // 4 lines (Tetris) = 1200 points at level 1
	}

	for _, tc := range testCases {
		game.Score = 0
		game.LinesCleared = 0

		game.addScore(tc.linesCleared)

		if game.Score != tc.expectedScore {
			t.Errorf("Expected score %d for %d lines, got %d", tc.expectedScore, tc.linesCleared, game.Score)
		}

		if game.LinesCleared != tc.linesCleared {
			t.Errorf("Expected lines cleared %d, got %d", tc.linesCleared, game.LinesCleared)
		}
	}

	// Test level up
	game.Score = 0
	game.LinesCleared = 9
	game.Level = 1

	game.addScore(1) // 10 lines total, should level up

	if game.Level != 2 {
		t.Errorf("Expected level to increase to 2, got %d", game.Level)
	}
}

func TestGameOver(t *testing.T) {
	game := NewGame()
	game.Start()

	// Create a situation where the game will end
	// Force a piece that can't be placed
	game.CurrentPiece = NewPiece(TypeI)

	// Fill the top rows to block piece placement
	for y := 0; y < 1; y++ {
		for x := 0; x < BoardWidth; x++ {
			game.Board.Cells[y][x] = I
		}
	}

	// Try to hold the piece (which should fail and trigger game over)
	game.HoldPiece()

	if game.State != StateGameOver {
		t.Errorf("Expected game state to be StateGameOver, got %d", game.State)
	}
}
