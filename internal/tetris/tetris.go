package tetris

// Package tetris implements a Tetris game engine.
// It provides the core game logic, board representation,
// piece movement, and scoring mechanics.

// GetColorForCell returns the color index for a given cell type
func GetColorForCell(cell Cell) int {
	return int(cell)
}

// GetPreviewPiece returns a copy of the next piece for display
func (g *Game) GetPreviewPiece() *Piece {
	return g.NextPiece.Copy()
}

// GetLevel returns the current game level
func (g *Game) GetLevel() int {
	return g.Level
}

// GetScore returns the current game score
func (g *Game) GetScore() int {
	return g.Score
}

// GetLinesCleared returns the total number of lines cleared
func (g *Game) GetLinesCleared() int {
	return g.LinesCleared
}

// IsGameOver returns true if the game is over
func (g *Game) IsGameOver() bool {
	return g.State == StateGameOver
}

// IsPaused returns true if the game is paused
func (g *Game) IsPaused() bool {
	return g.State == StatePaused
}

// IsPlaying returns true if the game is in progress
func (g *Game) IsPlaying() bool {
	return g.State == StatePlaying
}

// IsInMenu returns true if the game is in the menu
func (g *Game) IsInMenu() bool {
	return g.State == StateMenu
}
// GetHeldPiece returns a copy of the held piece for display
func (g *Game) GetHeldPiece() *Piece {
	if g.HeldPiece == nil {
		return nil
	}
	return g.HeldPiece.Copy()
}
