package tetris

import (
	"testing"

	_ "github.com/briancain/go-tetris/internal/testutil" // Import for init side effects
)

func TestGetTetriminoName(t *testing.T) {
	// Test that each piece type returns the correct name
	testCases := []struct {
		pieceType PieceType
		expected  string
	}{
		{TypeI, TetriminoI},
		{TypeJ, TetriminoJ},
		{TypeL, TetriminoL},
		{TypeO, TetriminoO},
		{TypeS, TetriminoS},
		{TypeT, TetriminoT},
		{TypeZ, TetriminoZ},
		{PieceType(99), "Unknown"}, // Invalid type should return "Unknown"
	}

	for _, tc := range testCases {
		result := GetTetriminoName(tc.pieceType)
		if result != tc.expected {
			t.Errorf("GetTetriminoName(%d) = %s, expected %s", tc.pieceType, result, tc.expected)
		}
	}
}

func TestTetriminoConstants(t *testing.T) {
	// Verify that the constants match the Tetris Guidelines
	if TetriminoI != "I-Tetrimino" {
		t.Errorf("Expected I-Tetrimino, got %s", TetriminoI)
	}
	if TetriminoJ != "J-Tetrimino" {
		t.Errorf("Expected J-Tetrimino, got %s", TetriminoJ)
	}
	if TetriminoL != "L-Tetrimino" {
		t.Errorf("Expected L-Tetrimino, got %s", TetriminoL)
	}
	if TetriminoO != "O-Tetrimino" {
		t.Errorf("Expected O-Tetrimino, got %s", TetriminoO)
	}
	if TetriminoS != "S-Tetrimino" {
		t.Errorf("Expected S-Tetrimino, got %s", TetriminoS)
	}
	if TetriminoT != "T-Tetrimino" {
		t.Errorf("Expected T-Tetrimino, got %s", TetriminoT)
	}
	if TetriminoZ != "Z-Tetrimino" {
		t.Errorf("Expected Z-Tetrimino, got %s", TetriminoZ)
	}
}
