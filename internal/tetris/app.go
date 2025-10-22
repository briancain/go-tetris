package tetris

import (
	"unicode"

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
	// Always update game state first (for multiplayer message processing)
	g.game.Update()

	// Handle input based on game state
	switch g.game.State {
	case StateMainMenu:
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			// Single Player
			g.game.Start()
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			// Multiplayer
			g.game.State = StateMultiplayerSetup
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			// Quit (for now, just do nothing)
		}
	case StateMultiplayerSetup:
		// Handle text input for username
		g.handleTextInput()

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			// Start multiplayer connection
			go g.game.StartMultiplayerConnection()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			// Back to main menu
			g.game.State = StateMainMenu
			g.game.UsernameInput = ""
			g.game.ConnectionStatus = ""
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			g.game.RemoveFromUsernameInput()
		}
	case StateMatchmaking:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			// Cancel matchmaking, back to main menu
			g.game.State = StateMainMenu
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

// handleTextInput processes text input for username entry
func (g *App) handleTextInput() {
	// Get typed characters
	runes := ebiten.AppendInputChars(nil)
	for _, r := range runes {
		// Only allow alphanumeric characters and basic symbols
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			g.game.AddToUsernameInput(r)
		}
	}
}
