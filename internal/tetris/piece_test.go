package tetris

import (
	"testing"

	_ "github.com/briancain/go-tetris/internal/testutil" // Import for init side effects
)

func TestNewPiece(t *testing.T) {
	// Test each piece type
	for pieceType := TypeI; pieceType <= TypeZ; pieceType++ {
		piece := NewPiece(pieceType)

		if piece == nil {
			t.Fatalf("NewPiece(%d) returned nil", pieceType)
		}

		if piece.Type != pieceType {
			t.Errorf("Expected piece type to be %d, got %d", pieceType, piece.Type)
		}

		// Check that the shape is not empty
		hasBlock := false
		for i := 0; i < len(piece.Shape); i++ {
			for j := 0; j < len(piece.Shape[i]); j++ {
				if piece.Shape[i][j] {
					hasBlock = true
					break
				}
			}
			if hasBlock {
				break
			}
		}

		if !hasBlock {
			t.Errorf("Piece type %d has an empty shape", pieceType)
		}
	}
}

func TestRandomPiece(t *testing.T) {
	// Test that RandomPiece returns a valid piece
	for i := 0; i < 100; i++ {
		piece := RandomPiece()

		if piece == nil {
			t.Fatal("RandomPiece returned nil")
		}

		if piece.Type < TypeI || piece.Type > TypeZ {
			t.Errorf("RandomPiece returned invalid piece type: %d", piece.Type)
		}
	}
}

func TestRotate(t *testing.T) {
	// Test rotation of T piece (I piece might not change visibly after rotation)
	piece := NewPiece(TypeT)
	originalShape := make([][]bool, len(piece.Shape))
	for i := range piece.Shape {
		originalShape[i] = make([]bool, len(piece.Shape[i]))
		copy(originalShape[i], piece.Shape[i])
	}

	// Rotate once
	piece.Rotate()

	// Check that the shape changed
	shapeChanged := false
	for i := 0; i < len(piece.Shape); i++ {
		for j := 0; j < len(piece.Shape[i]); j++ {
			if i < len(originalShape) && j < len(originalShape[i]) && piece.Shape[i][j] != originalShape[i][j] {
				shapeChanged = true
				break
			}
		}
		if shapeChanged {
			break
		}
	}

	if !shapeChanged {
		t.Error("Expected shape to change after rotation")
	}

	// Rotate three more times to get back to the original shape
	piece.Rotate()
	piece.Rotate()
	piece.Rotate()

	// Check that we're back to the original shape
	for i := 0; i < len(piece.Shape); i++ {
		for j := 0; j < len(piece.Shape[i]); j++ {
			if i < len(originalShape) && j < len(originalShape[i]) && piece.Shape[i][j] != originalShape[i][j] {
				t.Errorf("Expected to return to original shape after 4 rotations")
				return
			}
		}
	}
}

func TestCopy(t *testing.T) {
	original := NewPiece(TypeT)
	original.X = 5
	original.Y = 10

	pieceCopy := original.Copy()

	// Check that the copy has the same properties
	if pieceCopy.Type != original.Type {
		t.Errorf("Expected copy to have type %d, got %d", original.Type, pieceCopy.Type)
	}

	if pieceCopy.X != original.X {
		t.Errorf("Expected copy to have X %d, got %d", original.X, pieceCopy.X)
	}

	if pieceCopy.Y != original.Y {
		t.Errorf("Expected copy to have Y %d, got %d", original.Y, pieceCopy.Y)
	}

	// Check that the shape is the same
	for i := 0; i < len(original.Shape); i++ {
		for j := 0; j < len(original.Shape[i]); j++ {
			if pieceCopy.Shape[i][j] != original.Shape[i][j] {
				t.Errorf("Shape mismatch at (%d,%d): expected %v, got %v",
					j, i, original.Shape[i][j], pieceCopy.Shape[i][j])
			}
		}
	}

	// Modify the copy and check that the original is unchanged
	pieceCopy.X = 20
	pieceCopy.Shape[0][0] = !pieceCopy.Shape[0][0]

	if original.X != 5 {
		t.Error("Original X was modified when copy was changed")
	}

	if original.Shape[0][0] == pieceCopy.Shape[0][0] {
		t.Error("Original shape was modified when copy was changed")
	}
}
