package tetris

import (
	"time"
)

// Game states
const (
	StateMenu = iota
	StatePlaying
	StatePaused
	StateGameOver
)

// Game represents the Tetris game state
type Game struct {
	Board         *Board
	CurrentPiece  *Piece
	NextPiece     *Piece
	HeldPiece     *Piece // Piece that is being held
	HasSwapped    bool   // Flag to prevent multiple swaps per turn
	State         int
	Score         int
	Level         int
	LinesCleared  int
	DropTimer     time.Time
	DropInterval  time.Duration
	LastMoveDown  time.Time
	LastMoveSide  time.Time
	LastRotate    time.Time
	LastHold      time.Time // Time of last hold action
	InputDelay    time.Duration
	FastDropDelay time.Duration
}

// NewGame creates a new Tetris game
func NewGame() *Game {
	game := &Game{
		Board:         NewBoard(),
		State:         StateMenu,
		Score:         0,
		Level:         1,
		LinesCleared:  0,
		DropInterval:  800 * time.Millisecond, // Initial drop speed
		InputDelay:    100 * time.Millisecond, // Delay between input actions
		FastDropDelay: 50 * time.Millisecond,  // Fast drop speed
	}

	// Initialize pieces
	game.NextPiece = RandomPiece()

	return game
}

// Start begins a new game
func (g *Game) Start() {
	g.Board.Clear()
	g.CurrentPiece = g.NextPiece
	g.NextPiece = RandomPiece()
	g.HeldPiece = nil
	g.HasSwapped = false
	g.Score = 0
	g.Level = 1
	g.LinesCleared = 0
	g.State = StatePlaying
	g.DropTimer = time.Now()
	g.updateDropInterval()
}

// Update updates the game state
func (g *Game) Update() {
	if g.State != StatePlaying {
		return
	}

	// Check if it's time to drop the piece
	if time.Since(g.DropTimer) >= g.DropInterval {
		g.DropTimer = time.Now()
		g.moveDown()
	}
}

// moveDown moves the current piece down one row
func (g *Game) moveDown() {
	// Create a copy of the current piece and try to move it down
	testPiece := g.CurrentPiece.Copy()
	testPiece.Move(0, 1)

	// Check if the move is valid
	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Move(0, 1)
		return
	}

	// If we can't move down, lock the piece in place
	g.lockPiece()

	// Check for completed lines
	linesCleared := g.Board.ClearLines()
	if linesCleared > 0 {
		g.addScore(linesCleared)
	}

	// Check for game over
	g.CurrentPiece = g.NextPiece
	g.NextPiece = RandomPiece()

	if !g.Board.IsValidPosition(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y) {
		g.State = StateGameOver
	}
}

// MoveLeft moves the current piece left
func (g *Game) MoveLeft() bool {
	if g.State != StatePlaying || time.Since(g.LastMoveSide) < g.InputDelay {
		return false
	}

	g.LastMoveSide = time.Now()
	testPiece := g.CurrentPiece.Copy()
	testPiece.Move(-1, 0)

	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Move(-1, 0)
		return true
	}

	return false
}

// MoveRight moves the current piece right
func (g *Game) MoveRight() bool {
	if g.State != StatePlaying || time.Since(g.LastMoveSide) < g.InputDelay {
		return false
	}

	g.LastMoveSide = time.Now()
	testPiece := g.CurrentPiece.Copy()
	testPiece.Move(1, 0)

	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Move(1, 0)
		return true
	}

	return false
}

// RotatePiece rotates the current piece
func (g *Game) RotatePiece() bool {
	if g.State != StatePlaying || time.Since(g.LastRotate) < g.InputDelay {
		return false
	}

	g.LastRotate = time.Now()
	testPiece := g.CurrentPiece.Copy()
	testPiece.Rotate()

	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Rotate()
		return true
	}

	// Wall kick - try to adjust position if rotation fails
	// Try moving left
	testPiece.Move(-1, 0)
	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Rotate()
		g.CurrentPiece.Move(-1, 0)
		return true
	}

	// Try moving right
	testPiece.Move(2, 0) // Move 2 to the right from the left position
	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Rotate()
		g.CurrentPiece.Move(1, 0)
		return true
	}

	return false
}

// HardDrop drops the piece all the way down
func (g *Game) HardDrop() {
	if g.State != StatePlaying {
		return
	}

	// Keep moving down until we hit something
	for {
		testPiece := g.CurrentPiece.Copy()
		testPiece.Move(0, 1)

		if !g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
			break
		}

		g.CurrentPiece.Move(0, 1)
		g.Score++ // Small bonus for hard drop
	}

	g.lockPiece()

	// Check for completed lines
	linesCleared := g.Board.ClearLines()
	if linesCleared > 0 {
		g.addScore(linesCleared)
	}

	// Check for game over
	g.CurrentPiece = g.NextPiece
	g.NextPiece = RandomPiece()

	if !g.Board.IsValidPosition(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y) {
		g.State = StateGameOver
	}
}

// SoftDrop accelerates the piece downward
func (g *Game) SoftDrop() bool {
	if g.State != StatePlaying || time.Since(g.LastMoveDown) < g.FastDropDelay {
		return false
	}

	g.LastMoveDown = time.Now()
	testPiece := g.CurrentPiece.Copy()
	testPiece.Move(0, 1)

	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Move(0, 1)
		g.Score++ // Small bonus for soft drop
		return true
	}

	return false
}

// TogglePause toggles the game's pause state
func (g *Game) TogglePause() {
	if g.State == StatePlaying {
		g.State = StatePaused
	} else if g.State == StatePaused {
		g.State = StatePlaying
		g.DropTimer = time.Now() // Reset drop timer when unpausing
	}
}

// lockPiece locks the current piece in place on the board
func (g *Game) lockPiece() {
	g.Board.PlacePiece(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y, true)
	g.HasSwapped = false // Reset swap flag when piece is locked
}

// addScore adds to the score based on lines cleared
func (g *Game) addScore(linesCleared int) {
	// Classic Tetris scoring
	linePoints := []int{0, 40, 100, 300, 1200}
	g.Score += linePoints[linesCleared] * g.Level

	g.LinesCleared += linesCleared

	// Level up every 10 lines
	newLevel := (g.LinesCleared / 10) + 1
	if newLevel > g.Level {
		g.Level = newLevel
		g.updateDropInterval()
	}
}

// updateDropInterval adjusts the piece drop speed based on level
func (g *Game) updateDropInterval() {
	// Formula: 800ms * (0.8^(level-1))
	// This makes the game get progressively faster with each level
	baseInterval := 800.0
	factor := 0.8

	// Calculate the new interval
	interval := baseInterval * pow(factor, float64(g.Level-1))
	g.DropInterval = time.Duration(interval) * time.Millisecond
}

// pow calculates x^y for our drop interval calculation
func pow(x, y float64) float64 {
	result := 1.0
	for i := 0; i < int(y); i++ {
		result *= x
	}
	return result
}

// HoldPiece swaps the current piece with the held piece
func (g *Game) HoldPiece() bool {
	if g.State != StatePlaying || time.Since(g.LastHold) < g.InputDelay {
		return false
	}

	// Can only swap once per piece
	if g.HasSwapped {
		return false
	}

	g.LastHold = time.Now()

	// If there's no held piece yet, store current piece and get next piece
	if g.HeldPiece == nil {
		g.HeldPiece = g.CurrentPiece.Copy()
		// Reset the held piece to its original position and orientation
		g.HeldPiece = NewPiece(g.HeldPiece.Type)

		g.CurrentPiece = g.NextPiece
		g.NextPiece = RandomPiece()
	} else {
		// Swap current piece with held piece
		tempPiece := g.CurrentPiece

		// Create a new piece of the held type at the top of the board
		g.CurrentPiece = NewPiece(g.HeldPiece.Type)

		// Store the previous current piece as held
		g.HeldPiece = NewPiece(tempPiece.Type)
	}

	// Mark that we've swapped this turn
	g.HasSwapped = true

	// Check if the new current piece can be placed
	if !g.Board.IsValidPosition(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y) {
		// Game over if the piece can't be placed
		g.State = StateGameOver
		return false
	}

	return true
}
