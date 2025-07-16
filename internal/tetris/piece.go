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

// Rotation states
const (
	RotationState0 = iota // Initial orientation
	RotationState1        // 90 degrees clockwise
	RotationState2        // 180 degrees
	RotationState3        // 270 degrees clockwise
)

// SRS wall kick data for J, L, S, T, Z pieces
var wallKickDataJLSTZ = [][][]int{
	{ // 0->1
		{-1, 0}, {-1, -1}, {0, 2}, {-1, 2},
	},
	{ // 1->2
		{1, 0}, {1, 1}, {0, -2}, {1, -2},
	},
	{ // 2->3
		{1, 0}, {1, -1}, {0, 2}, {1, 2},
	},
	{ // 3->0
		{-1, 0}, {-1, 1}, {0, -2}, {-1, -2},
	},
}

// SRS wall kick data for I piece
var wallKickDataI = [][][]int{
	{ // 0->1
		{-2, 0}, {1, 0}, {-2, 1}, {1, -2},
	},
	{ // 1->2
		{-1, 0}, {2, 0}, {-1, -2}, {2, 1},
	},
	{ // 2->3
		{2, 0}, {-1, 0}, {2, -1}, {-1, 2},
	},
	{ // 3->0
		{1, 0}, {-2, 0}, {1, 2}, {-2, -1},
	},
}

// Piece represents a tetromino piece
type Piece struct {
	Type          PieceType
	Shape         [][]bool
	X             int
	Y             int
	RotationState int // Current rotation state (0-3)
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
	var x int
	if pieceType == TypeI || pieceType == TypeO {
		// I and O spawn in the middle columns
		x = (BoardWidth - len(shape[0])) / 2
	} else {
		// Others spawn in the left-middle columns
		x = (BoardWidth-len(shape[0]))/2 - 1
	}
	y := 0

	return &Piece{
		Type:          pieceType,
		Shape:         shapeCopy,
		X:             x,
		Y:             y,
		RotationState: RotationState0,
	}
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

	// Update rotation state
	p.RotationState = (p.RotationState + 1) % 4
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
		Type:          p.Type,
		Shape:         shapeCopy,
		X:             p.X,
		Y:             p.Y,
		RotationState: p.RotationState,
	}
}

// PieceGenerator manages the random generation of pieces using the 7-bag system
type PieceGenerator struct {
	bag      []PieceType
	bagIndex int
	rng      *rand.Rand
}

// NewPieceGenerator creates a new piece generator with the 7-bag system
func NewPieceGenerator() *PieceGenerator {
	generator := &PieceGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	generator.refillBag()
	return generator
}

// refillBag creates a new shuffled bag of all 7 piece types
func (pg *PieceGenerator) refillBag() {
	// Create a bag with all 7 piece types
	pg.bag = []PieceType{TypeI, TypeJ, TypeL, TypeO, TypeS, TypeT, TypeZ}

	// Shuffle the bag
	pg.rng.Shuffle(len(pg.bag), func(i, j int) {
		pg.bag[i], pg.bag[j] = pg.bag[j], pg.bag[i]
	})

	pg.bagIndex = 0
}

// NextPiece gets the next piece from the bag
func (pg *PieceGenerator) NextPiece() *Piece {
	// If we've used all pieces in the bag, refill it
	if pg.bagIndex >= len(pg.bag) {
		pg.refillBag()
	}

	// Get the next piece type from the bag
	pieceType := pg.bag[pg.bagIndex]
	pg.bagIndex++

	// Create and return a new piece of this type
	return NewPiece(pieceType)
}
