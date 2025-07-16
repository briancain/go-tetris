package tetris

import (
	"testing"

	_ "github.com/briancain/go-tetris/internal/testutil" // Import for init side effects
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	if board == nil {
		t.Fatal("NewBoard returned nil")
	}

	// Check board dimensions
	if len(board.Cells) != BoardHeight {
		t.Errorf("Expected board height to be %d, got %d", BoardHeight, len(board.Cells))
	}

	for i := 0; i < BoardHeight; i++ {
		if len(board.Cells[i]) != BoardWidth {
			t.Errorf("Expected board width at row %d to be %d, got %d", i, BoardWidth, len(board.Cells[i]))
		}
	}

	// Check that all cells are empty
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if board.Cells[y][x] != Empty {
				t.Errorf("Expected cell at (%d,%d) to be Empty, got %d", x, y, board.Cells[y][x])
			}
		}
	}
}

func TestIsValidPosition(t *testing.T) {
	board := NewBoard()
	piece := NewPiece(TypeI)

	// Valid position
	if !board.IsValidPosition(piece, 3, 0) {
		t.Error("Expected position (3,0) to be valid")
	}

	// Invalid position - out of bounds horizontally
	if board.IsValidPosition(piece, -1, 0) {
		t.Error("Expected position (-1,0) to be invalid")
	}

	if board.IsValidPosition(piece, BoardWidth, 0) {
		t.Error("Expected position (BoardWidth,0) to be invalid")
	}

	// Invalid position - out of bounds vertically
	if board.IsValidPosition(piece, 3, BoardHeight) {
		t.Error("Expected position (3,BoardHeight) to be invalid")
	}

	// Test collision with existing blocks
	board.Cells[5][3] = I
	piece.Y = 3
	if board.IsValidPosition(piece, 3, 5) {
		t.Error("Expected position to be invalid due to collision")
	}
}

func TestPlacePiece(t *testing.T) {
	board := NewBoard()
	piece := NewPiece(TypeI)
	piece.X = 3
	piece.Y = 0

	board.PlacePiece(piece, piece.X, piece.Y, false)

	// Check that the piece cells are now on the board
	for i := 0; i < len(piece.Shape); i++ {
		for j := 0; j < len(piece.Shape[i]); j++ {
			if piece.Shape[i][j] {
				if Cell(piece.Type) != board.Cells[piece.Y+i][piece.X+j] {
					t.Errorf("Expected cell at (%d,%d) to be %d, got %d",
						piece.X+j, piece.Y+i, Cell(piece.Type), board.Cells[piece.Y+i][piece.X+j])
				}
			}
		}
	}
}

func TestClearLines(t *testing.T) {
	board := NewBoard()

	// Fill a line
	for x := 0; x < BoardWidth; x++ {
		board.Cells[BoardHeight-1][x] = I
	}

	// Clear lines and check result
	linesCleared := board.ClearLines()

	if linesCleared != 1 {
		t.Errorf("Expected 1 line cleared, got %d", linesCleared)
	}

	// Check that the line is now empty
	for x := 0; x < BoardWidth; x++ {
		if board.Cells[BoardHeight-1][x] != Empty {
			t.Errorf("Expected cell at (%d,%d) to be Empty after clearing, got %d",
				x, BoardHeight-1, board.Cells[BoardHeight-1][x])
		}
	}

	// Test multiple lines
	for y := BoardHeight - 3; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			board.Cells[y][x] = I
		}
	}

	linesCleared = board.ClearLines()

	if linesCleared != 3 {
		t.Errorf("Expected 3 lines cleared, got %d", linesCleared)
	}
}

func TestIsLineFull(t *testing.T) {
	board := NewBoard()

	// Empty line
	if board.isLineFull(0) {
		t.Error("Expected empty line to not be full")
	}

	// Partially filled line
	board.Cells[0][0] = I
	if board.isLineFull(0) {
		t.Error("Expected partially filled line to not be full")
	}

	// Full line
	for x := 0; x < BoardWidth; x++ {
		board.Cells[1][x] = I
	}
	if !board.isLineFull(1) {
		t.Error("Expected full line to be full")
	}
}
