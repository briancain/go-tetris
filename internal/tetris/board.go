package tetris

// Board dimensions
const (
	BoardWidth  = 10
	BoardHeight = 20
)

// Cell represents a single cell in the game board
type Cell int

// Cell states
const (
	Empty Cell = iota
	I
	J
	L
	O
	S
	T
	Z
	Locked // For pieces that have been locked in place
)

// Board represents the Tetris game board
type Board struct {
	Cells [BoardHeight][BoardWidth]Cell
}

// NewBoard creates a new empty board
func NewBoard() *Board {
	board := &Board{}
	board.Clear()
	return board
}

// Clear resets the board to empty
func (b *Board) Clear() {
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			b.Cells[y][x] = Empty
		}
	}
}

// IsValidPosition checks if a piece can be placed at the given position
func (b *Board) IsValidPosition(piece *Piece, x, y int) bool {
	for i := 0; i < len(piece.Shape); i++ {
		for j := 0; j < len(piece.Shape[i]); j++ {
			if piece.Shape[i][j] {
				newX, newY := x+j, y+i

				// Check if out of bounds
				if newX < 0 || newX >= BoardWidth || newY < 0 || newY >= BoardHeight {
					return false
				}

				// Check if cell is already occupied
				if b.Cells[newY][newX] != Empty {
					return false
				}
			}
		}
	}
	return true
}

// PlacePiece places a piece on the board
func (b *Board) PlacePiece(piece *Piece, x, y int, _ bool) {
	// Always use the piece's original color type
	cellType := Cell(piece.Type)

	for i := 0; i < len(piece.Shape); i++ {
		for j := 0; j < len(piece.Shape[i]); j++ {
			if piece.Shape[i][j] {
				b.Cells[y+i][x+j] = cellType
			}
		}
	}
}

// ClearLines removes completed lines and returns the number of lines cleared
func (b *Board) ClearLines() int {
	linesCleared := 0

	for y := BoardHeight - 1; y >= 0; y-- {
		if b.isLineFull(y) {
			b.clearLine(y)
			linesCleared++
			y++ // Check the same line again after shifting
		}
	}

	return linesCleared
}

// isLineFull checks if a line is completely filled
func (b *Board) isLineFull(y int) bool {
	for x := 0; x < BoardWidth; x++ {
		if b.Cells[y][x] == Empty {
			return false
		}
	}
	return true
}

// clearLine removes a line and shifts all lines above down
func (b *Board) clearLine(y int) {
	// Move all lines above down
	for i := y; i > 0; i-- {
		for x := 0; x < BoardWidth; x++ {
			b.Cells[i][x] = b.Cells[i-1][x]
		}
	}

	// Clear the top line
	for x := 0; x < BoardWidth; x++ {
		b.Cells[0][x] = Empty
	}
}
