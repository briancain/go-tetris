package tetris

import (
	"math/rand"
	"time"
)

// PieceType represents the type of tetromino
type PieceType int

// Tetromino types
const (
	TypeI PieceType = iota + 1
	TypeJ
	TypeL
	TypeO
	TypeS
	TypeT
	TypeZ
)

// Piece represents a tetromino piece
type Piece struct {
	Type  PieceType
	Shape [][]bool
	X     int
	Y     int
}

// Tetromino shapes
var (
	// I-piece: ####
	shapeI = [][]bool{
		{true, true, true, true},
	}

	// J-piece: #
	//          ###
	shapeJ = [][]bool{
		{true, false, false},
		{true, true, true},
	}

	// L-piece:   #
	//          ###
	shapeL = [][]bool{
		{false, false, true},
		{true, true, true},
	}

	// O-piece: ##
	//          ##
	shapeO = [][]bool{
		{true, true},
		{true, true},
	}

	// S-piece:  ##
	//          ##
	shapeS = [][]bool{
		{false, true, true},
		{true, true, false},
	}

	// T-piece:  #
	//          ###
	shapeT = [][]bool{
		{false, true, false},
		{true, true, true},
	}

	// Z-piece: ##
	//           ##
	shapeZ = [][]bool{
		{true, true, false},
		{false, true, true},
	}
)

// NewPiece creates a new piece of the specified type
func NewPiece(pieceType PieceType) *Piece {
	var shape [][]bool

	switch pieceType {
	case TypeI:
		shape = shapeI
	case TypeJ:
		shape = shapeJ
	case TypeL:
		shape = shapeL
	case TypeO:
		shape = shapeO
	case TypeS:
		shape = shapeS
	case TypeT:
		shape = shapeT
	case TypeZ:
		shape = shapeZ
	}

	// Deep copy the shape to avoid modifying the original
	shapeCopy := make([][]bool, len(shape))
	for i := range shape {
		shapeCopy[i] = make([]bool, len(shape[i]))
		copy(shapeCopy[i], shape[i])
	}

	// Start position at the top center of the board
	x := (BoardWidth - len(shape[0])) / 2
	y := 0

	return &Piece{
		Type:  pieceType,
		Shape: shapeCopy,
		X:     x,
		Y:     y,
	}
}

// RandomPiece creates a random tetromino piece
func RandomPiece() *Piece {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pieceType := PieceType(r.Intn(7) + 1) // 1-7 for the different piece types
	return NewPiece(pieceType)
}

// Rotate rotates the piece clockwise
func (p *Piece) Rotate() {
	// Skip rotation for O piece (square)
	if p.Type == TypeO {
		return
	}

	// Get dimensions
	rows := len(p.Shape)
	cols := len(p.Shape[0])

	// Create a new rotated shape
	rotated := make([][]bool, cols)
	for i := range rotated {
		rotated[i] = make([]bool, rows)
	}

	// Perform rotation
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			rotated[j][rows-1-i] = p.Shape[i][j]
		}
	}

	p.Shape = rotated
}

// Move moves the piece by the given delta
func (p *Piece) Move(dx, dy int) {
	p.X += dx
	p.Y += dy
}

// Copy creates a deep copy of the piece
func (p *Piece) Copy() *Piece {
	shapeCopy := make([][]bool, len(p.Shape))
	for i := range p.Shape {
		shapeCopy[i] = make([]bool, len(p.Shape[i]))
		copy(shapeCopy[i], p.Shape[i])
	}

	return &Piece{
		Type:  p.Type,
		Shape: shapeCopy,
		X:     p.X,
		Y:     p.Y,
	}
}
