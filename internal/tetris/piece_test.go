package tetris

import (
	"testing"
)

func TestNewPiece(t *testing.T) {
	testCases := []struct {
		pieceType PieceType
		width     int
		height    int
	}{
		{TypeI, 4, 1}, // I piece is 4x1
		{TypeJ, 3, 2}, // J piece is 3x2
		{TypeL, 3, 2}, // L piece is 3x2
		{TypeO, 2, 2}, // O piece is 2x2
		{TypeS, 3, 2}, // S piece is 3x2
		{TypeT, 3, 2}, // T piece is 3x2
		{TypeZ, 3, 2}, // Z piece is 3x2
	}
	
	for _, tc := range testCases {
		piece := NewPiece(tc.pieceType)
		
		// Check piece type
		if piece.Type != tc.pieceType {
			t.Errorf("Expected piece type %v, got %v", tc.pieceType, piece.Type)
		}
		
		// Check shape dimensions
		if len(piece.Shape) != tc.height {
			t.Errorf("Expected piece height %d, got %d", tc.height, len(piece.Shape))
		}
		
		if len(piece.Shape[0]) != tc.width {
			t.Errorf("Expected piece width %d, got %d", tc.width, len(piece.Shape[0]))
		}
		
		// Check that the piece is positioned at the top center of the board
		expectedX := (BoardWidth - len(piece.Shape[0])) / 2
		if piece.X != expectedX {
			t.Errorf("Expected piece X position %d, got %d", expectedX, piece.X)
		}
		
		if piece.Y != 0 {
			t.Errorf("Expected piece Y position 0, got %d", piece.Y)
		}
	}
}

func TestRandomPiece(t *testing.T) {
	// Test that RandomPiece returns a valid piece
	piece := RandomPiece()
	
	if piece == nil {
		t.Error("RandomPiece returned nil")
	}
	
	// Check that the piece type is valid
	if piece.Type < TypeI || piece.Type > TypeZ {
		t.Errorf("Invalid piece type: %v", piece.Type)
	}
}

func TestRotate(t *testing.T) {
	// Test rotation of I piece
	piece := NewPiece(TypeI)
	originalWidth := len(piece.Shape[0])
	originalHeight := len(piece.Shape)
	
	piece.Rotate()
	
	// After rotation, width and height should be swapped
	if len(piece.Shape) != originalWidth {
		t.Errorf("Expected rotated height %d, got %d", originalWidth, len(piece.Shape))
	}
	
	if len(piece.Shape[0]) != originalHeight {
		t.Errorf("Expected rotated width %d, got %d", originalHeight, len(piece.Shape[0]))
	}
	
	// Test that O piece doesn't change when rotated
	piece = NewPiece(TypeO)
	originalShape := make([][]bool, len(piece.Shape))
	for i := range piece.Shape {
		originalShape[i] = make([]bool, len(piece.Shape[i]))
		copy(originalShape[i], piece.Shape[i])
	}
	
	piece.Rotate()
	
	// O piece should remain unchanged
	for i := range piece.Shape {
		for j := range piece.Shape[i] {
			if piece.Shape[i][j] != originalShape[i][j] {
				t.Errorf("O piece shape changed after rotation at (%d,%d)", j, i)
			}
		}
	}
}

func TestMove(t *testing.T) {
	piece := NewPiece(TypeT)
	originalX := piece.X
	originalY := piece.Y
	
	// Move right
	piece.Move(1, 0)
	if piece.X != originalX+1 {
		t.Errorf("Expected X position %d, got %d", originalX+1, piece.X)
	}
	if piece.Y != originalY {
		t.Errorf("Expected Y position %d, got %d", originalY, piece.Y)
	}
	
	// Move down
	piece.Move(0, 1)
	if piece.X != originalX+1 {
		t.Errorf("Expected X position %d, got %d", originalX+1, piece.X)
	}
	if piece.Y != originalY+1 {
		t.Errorf("Expected Y position %d, got %d", originalY+1, piece.Y)
	}
	
	// Move left and up
	piece.Move(-2, -2)
	if piece.X != originalX-1 {
		t.Errorf("Expected X position %d, got %d", originalX-1, piece.X)
	}
	if piece.Y != originalY-1 {
		t.Errorf("Expected Y position %d, got %d", originalY-1, piece.Y)
	}
}

func TestCopy(t *testing.T) {
	original := NewPiece(TypeT)
	original.X = 5
	original.Y = 10
	
	// Create a copy
	copy := original.Copy()
	
	// Check that the copy has the same properties
	if copy.Type != original.Type {
		t.Errorf("Expected piece type %v, got %v", original.Type, copy.Type)
	}
	
	if copy.X != original.X {
		t.Errorf("Expected X position %d, got %d", original.X, copy.X)
	}
	
	if copy.Y != original.Y {
		t.Errorf("Expected Y position %d, got %d", original.Y, copy.Y)
	}
	
	// Check that the shape is a deep copy
	if len(copy.Shape) != len(original.Shape) {
		t.Errorf("Expected shape height %d, got %d", len(original.Shape), len(copy.Shape))
	}
	
	for i := range original.Shape {
		if len(copy.Shape[i]) != len(original.Shape[i]) {
			t.Errorf("Expected shape width %d at row %d, got %d", len(original.Shape[i]), i, len(copy.Shape[i]))
		}
		
		for j := range original.Shape[i] {
			if copy.Shape[i][j] != original.Shape[i][j] {
				t.Errorf("Shape mismatch at (%d,%d)", j, i)
			}
		}
	}
	
	// Modify the copy and check that the original is unchanged
	copy.Shape[0][0] = !copy.Shape[0][0]
	if copy.Shape[0][0] == original.Shape[0][0] {
		t.Error("Modifying copy affected the original")
	}
}
