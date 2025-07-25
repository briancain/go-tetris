package tetris

// Board dimensions
const (
	BoardWidth            = 10
	BoardHeight           = 20
	BoardHeightWithBuffer = 22 // Total height including hidden rows
)

// Cell represents a single cell in the game board
type Cell int

// Cell states
const (
	Empty   Cell = iota
	CyanI        // I - Cyan
	BlueJ        // J - Blue
	OrangeL      // L - Orange
	YellowO      // O - Yellow
	GreenS       // S - Green
	PurpleT      // T - Purple
	RedZ         // Z - Red
	Locked       // For pieces that have been locked in place
)

// Board represents the Tetris game board
type Board struct {
	Cells [BoardHeightWithBuffer][BoardWidth]Cell
}

// NewBoard creates a new empty board
func NewBoard() *Board {
	board := &Board{}
	board.Clear()
	return board
}

// Clear resets the board to empty
func (b *Board) Clear() {
	for y := 0; y < BoardHeightWithBuffer; y++ {
		for x := 0; x < BoardWidth; x++ {
			b.Cells[y][x] = Empty
		}
	}
}

// IsValidPosition checks if a piece can be placed at the given position
func (b *Board) IsValidPosition(piece *Piece, x, y int) bool {
	shape := piece.Shape
	height := len(shape)
	width := len(shape[0])

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if !shape[i][j] {
				continue // Skip empty cells in the piece shape
			}

			newX, newY := x+j, y+i

			// Check if out of bounds
			if newX < 0 || newX >= BoardWidth || newY < 0 || newY >= BoardHeightWithBuffer {
				return false
			}

			// Check if cell is already occupied
			if b.Cells[newY][newX] != Empty {
				return false
			}
		}
	}
	return true
}

// PlacePiece places a piece on the board
func (b *Board) PlacePiece(piece *Piece, x, y int, _ bool) {
	// Map piece type to the correct cell type with guideline colors
	var cellType Cell
	switch piece.Type {
	case TypeI:
		cellType = CyanI
	case TypeJ:
		cellType = BlueJ
	case TypeL:
		cellType = OrangeL
	case TypeO:
		cellType = YellowO
	case TypeS:
		cellType = GreenS
	case TypeT:
		cellType = PurpleT
	case TypeZ:
		cellType = RedZ
	}

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
	destRow := BoardHeightWithBuffer - 1

	// Process rows from bottom to top
	for srcRow := BoardHeightWithBuffer - 1; srcRow >= BoardHeightWithBuffer-BoardHeight; srcRow-- {
		// If the line is full, don't copy it (effectively removing it)
		if b.isLineFull(srcRow) {
			linesCleared++
			continue
		}

		// If we're not at the same row, copy the row down
		if srcRow != destRow {
			// Copy the row
			copy(b.Cells[destRow][:], b.Cells[srcRow][:])
		}

		destRow--
	}

	// Clear the top rows
	for y := 0; y <= destRow; y++ {
		for x := 0; x < BoardWidth; x++ {
			b.Cells[y][x] = Empty
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
