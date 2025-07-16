package tetris

// constants.go contains constants and terminology as defined by the Tetris Guidelines

// Tetrimino names according to Tetris Guidelines
const (
	TetriminoI = "I-Tetrimino" // Cyan I
	TetriminoJ = "J-Tetrimino" // Blue J
	TetriminoL = "L-Tetrimino" // Orange L
	TetriminoO = "O-Tetrimino" // Yellow O
	TetriminoS = "S-Tetrimino" // Green S
	TetriminoT = "T-Tetrimino" // Purple T
	TetriminoZ = "Z-Tetrimino" // Red Z
)

// Game mode names
const (
	ModeMarathon = "Marathon"
	ModeSprint   = "Sprint"
	ModeUltra    = "Ultra"
)

// Special move names
const (
	MoveTSpin        = "T-Spin"
	MoveTSpinMini    = "T-Spin Mini"
	MoveTetris       = "Tetris"
	MoveBackToBack   = "Back-to-Back"
	MovePerfectClear = "Perfect Clear"
)

// GetTetriminoName returns the official name for a piece type
func GetTetriminoName(pieceType PieceType) string {
	switch pieceType {
	case TypeI:
		return TetriminoI
	case TypeJ:
		return TetriminoJ
	case TypeL:
		return TetriminoL
	case TypeO:
		return TetriminoO
	case TypeS:
		return TetriminoS
	case TypeT:
		return TetriminoT
	case TypeZ:
		return TetriminoZ
	default:
		return "Unknown"
	}
}
