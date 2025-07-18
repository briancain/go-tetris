package tetris

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// App is the main game application
type App struct {
	game     *Game
	renderer interface {
		Draw(screen *ebiten.Image)
	}
}

// NewApp creates a new Tetris game application
func NewApp(game *Game, renderer interface {
	Draw(screen *ebiten.Image)
}) *App {
	return &App{
		game:     game,
		renderer: renderer,
	}
}

// Update updates the game state
func (g *App) Update() error {
	// Handle input based on game state
	switch g.game.State {
	case StateMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.game.Start()
		}
	case StatePlaying:
		// Game controls
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.game.TogglePause()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.game.HardDrop()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.game.RotatePiece()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) || inpututil.IsKeyJustPressed(ebiten.KeyShiftRight) {
			g.game.HoldPiece()
		}

		// Continuous movement
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.game.MoveLeft()
		}

		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.game.MoveRight()
		}

		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.game.SoftDrop()
		}

		// Update game state
		g.game.Update()
	case StatePaused:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.game.TogglePause()
		}
	case StateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.game.Start()
		}
	}

	return nil
}

// Draw draws the game
func (g *App) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen)
}

// Layout returns the game's logical screen size
func (g *App) Layout(_, _ int) (int, int) {
	return 640, 480 // Fixed game resolution
}
