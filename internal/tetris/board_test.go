package tetris

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()
	
	// Check that the board is initialized with empty cells
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if board.Cells[y][x] != Empty {
				t.Errorf("Expected empty cell at (%d,%d), got %v", x, y, board.Cells[y][x])
			}
		}
	}
}

func TestClear(t *testing.T) {
	board := NewBoard()
	
	// Fill some cells
	board.Cells[0][0] = I
	board.Cells[1][1] = J
	board.Cells[2][2] = L
	
	// Clear the board
	board.Clear()
	
	// Check that all cells are empty
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if board.Cells[y][x] != Empty {
				t.Errorf("Expected empty cell at (%d,%d), got %v", x, y, board.Cells[y][x])
			}
		}
	}
}

func TestIsValidPosition(t *testing.T) {
	board := NewBoard()
	piece := NewPiece(TypeI) // I piece is a horizontal line of 4 blocks
	
	// Test valid position
	if !board.IsValidPosition(piece, 3, 0) {
		t.Error("Expected position to be valid")
	}
	
	// Test position outside board boundaries (left)
	if board.IsValidPosition(piece, -1, 0) {
		t.Error("Expected position to be invalid (left boundary)")
	}
	
	// Test position outside board boundaries (right)
	if board.IsValidPosition(piece, BoardWidth-3, 0) {
		t.Error("Expected position to be invalid (right boundary)")
	}
	
	// Test position outside board boundaries (bottom)
	if board.IsValidPosition(piece, 3, BoardHeight) {
		t.Error("Expected position to be invalid (bottom boundary)")
	}
	
	// Test collision with existing piece
	board.Cells[5][5] = L
	piece = NewPiece(TypeO) // O piece is a 2x2 square
	if board.IsValidPosition(piece, 4, 4) {
		t.Error("Expected position to be invalid (collision)")
	}
}

func TestPlacePiece(t *testing.T) {
	board := NewBoard()
	piece := NewPiece(TypeT)
	
	// Place the piece
	board.PlacePiece(piece, 3, 0, false)
	
	// Check that the cells are filled correctly
	for i := 0; i < len(piece.Shape); i++ {
		for j := 0; j < len(piece.Shape[i]); j++ {
			if piece.Shape[i][j] {
				if board.Cells[0+i][3+j] != Cell(piece.Type) {
					t.Errorf("Expected cell at (%d,%d) to be %v, got %v", 3+j, 0+i, Cell(piece.Type), board.Cells[0+i][3+j])
				}
			}
		}
	}
}

func TestClearLines(t *testing.T) {
	board := NewBoard()
	
	// Fill a complete line
	for x := 0; x < BoardWidth; x++ {
		board.Cells[BoardHeight-1][x] = I
	}
	
	// Fill part of another line
	for x := 0; x < BoardWidth-1; x++ {
		board.Cells[BoardHeight-2][x] = J
	}
	
	// Clear lines and check the result
	linesCleared := board.ClearLines()
	
	if linesCleared != 1 {
		t.Errorf("Expected 1 line cleared, got %d", linesCleared)
	}
	
	// Check that the partial line was moved down
	for x := 0; x < BoardWidth-1; x++ {
		if board.Cells[BoardHeight-1][x] != J {
			t.Errorf("Expected cell at (%d,%d) to be J, got %v", x, BoardHeight-1, board.Cells[BoardHeight-1][x])
		}
	}
	
	// Check that the last cell of the moved line is empty
	if board.Cells[BoardHeight-1][BoardWidth-1] != Empty {
		t.Errorf("Expected empty cell at (%d,%d), got %v", BoardWidth-1, BoardHeight-1, board.Cells[BoardHeight-1][BoardWidth-1])
	}
	
	// Check that the top line is now empty
	for x := 0; x < BoardWidth; x++ {
		if board.Cells[BoardHeight-2][x] != Empty {
			t.Errorf("Expected empty cell at (%d,%d), got %v", x, BoardHeight-2, board.Cells[BoardHeight-2][x])
		}
	}
}

func TestIsLineFull(t *testing.T) {
	board := NewBoard()
	
	// Fill a complete line
	for x := 0; x < BoardWidth; x++ {
		board.Cells[BoardHeight-1][x] = I
	}
	
	// Fill part of another line
	for x := 0; x < BoardWidth-1; x++ {
		board.Cells[BoardHeight-2][x] = J
	}
	
	// Check if the complete line is detected as full
	if !board.isLineFull(BoardHeight-1) {
		t.Error("Expected line to be full")
	}
	
	// Check if the partial line is not detected as full
	if board.isLineFull(BoardHeight-2) {
		t.Error("Expected line to not be full")
	}
}
