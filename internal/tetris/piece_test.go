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

		// Check that rotation state is initialized to 0
		if piece.RotationState != RotationState0 {
			t.Errorf("Expected initial rotation state to be %d, got %d",
				RotationState0, piece.RotationState)
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

		// Check spawn position based on piece type
		if pieceType == TypeI || pieceType == TypeO {
			// I and O should spawn in the middle
			expectedX := (BoardWidth - len(piece.Shape[0])) / 2
			if piece.X != expectedX {
				t.Errorf("Expected I/O piece X to be %d, got %d", expectedX, piece.X)
			}
		} else {
			// Others should spawn in the left-middle
			expectedX := (BoardWidth-len(piece.Shape[0]))/2 - 1
			if piece.X != expectedX {
				t.Errorf("Expected piece X to be %d, got %d", expectedX, piece.X)
			}
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

	// Initial rotation state should be 0
	if piece.RotationState != RotationState0 {
		t.Errorf("Expected initial rotation state to be %d, got %d",
			RotationState0, piece.RotationState)
	}

	// Rotate once
	piece.Rotate()

	// Check that rotation state changed
	if piece.RotationState != RotationState1 {
		t.Errorf("Expected rotation state to be %d after one rotation, got %d",
			RotationState1, piece.RotationState)
	}

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
	piece.Rotate() // State 2
	if piece.RotationState != RotationState2 {
		t.Errorf("Expected rotation state to be %d after two rotations, got %d",
			RotationState2, piece.RotationState)
	}

	piece.Rotate() // State 3
	if piece.RotationState != RotationState3 {
		t.Errorf("Expected rotation state to be %d after three rotations, got %d",
			RotationState3, piece.RotationState)
	}

	piece.Rotate() // Back to State 0
	if piece.RotationState != RotationState0 {
		t.Errorf("Expected rotation state to be %d after four rotations, got %d",
			RotationState0, piece.RotationState)
	}

	// Check that we're back to the original shape
	for i := 0; i < len(piece.Shape); i++ {
		for j := 0; j < len(piece.Shape[i]); j++ {
			if i < len(originalShape) && j < len(originalShape[i]) && piece.Shape[i][j] != originalShape[i][j] {
				t.Errorf("Expected to return to original shape after 4 rotations")
				return
			}
		}
	}

	// Test O piece rotation (should not change)
	oPiece := NewPiece(TypeO)
	originalState := oPiece.RotationState
	oPiece.Rotate()

	if oPiece.RotationState != originalState {
		t.Errorf("O piece rotation state should not change, was %d, now %d",
			originalState, oPiece.RotationState)
	}
}

func TestCopy(t *testing.T) {
	original := NewPiece(TypeT)
	original.X = 5
	original.Y = 10
	original.RotationState = RotationState2 // Set a non-default rotation state

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

	if pieceCopy.RotationState != original.RotationState {
		t.Errorf("Expected copy to have rotation state %d, got %d",
			original.RotationState, pieceCopy.RotationState)
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
	pieceCopy.RotationState = RotationState3
	pieceCopy.Shape[0][0] = !pieceCopy.Shape[0][0]

	if original.X != 5 {
		t.Error("Original X was modified when copy was changed")
	}

	if original.RotationState != RotationState2 {
		t.Error("Original rotation state was modified when copy was changed")
	}

	if original.Shape[0][0] == pieceCopy.Shape[0][0] {
		t.Error("Original shape was modified when copy was changed")
	}
}
func TestPieceGenerator(t *testing.T) {
	generator := NewPieceGenerator()

	if generator == nil {
		t.Fatal("NewPieceGenerator returned nil")
	}

	// Test that the bag is initialized
	if len(generator.bag) != 7 {
		t.Errorf("Expected bag to have 7 pieces, got %d", len(generator.bag))
	}

	// Get 7 pieces and track which ones we've seen
	seenPieces := make(map[PieceType]bool)
	for i := 0; i < 7; i++ {
		piece := generator.NextPiece()
		seenPieces[piece.Type] = true
	}

	// Check that we got all 7 piece types
	expectedTypes := []PieceType{TypeI, TypeJ, TypeL, TypeO, TypeS, TypeT, TypeZ}
	for _, pieceType := range expectedTypes {
		if !seenPieces[pieceType] {
			t.Errorf("Expected to see piece type %d in the first bag", pieceType)
		}
	}

	// Check that the bag refills after 7 pieces
	piece8 := generator.NextPiece()
	if piece8 == nil {
		t.Error("Failed to get piece after bag should have refilled")
	}
}
