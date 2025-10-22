package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text" // nolint:staticcheck // Using deprecated API for compatibility
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"github.com/briancain/go-tetris/internal/tetris"
)

const (
	// Screen dimensions
	ScreenWidth  = 640
	ScreenHeight = 480

	// Cell size in pixels
	CellSize = 20

	// Board position on screen
	BoardX = 220
	BoardY = 40

	// Preview position on screen
	PreviewX = 480
	PreviewY = 80

	// Hold piece position on screen
	HoldX = 100
	HoldY = 80

	// UI colors
	BackgroundColor = 0x1A1A1AFF
)

// Tetris piece colors
var pieceColors = []color.RGBA{
	{0, 0, 0, 0},         // Empty
	{0, 255, 255, 255},   // CyanI - Cyan I
	{0, 0, 255, 255},     // BlueJ - Blue J
	{255, 165, 0, 255},   // OrangeL - Orange L
	{255, 255, 0, 255},   // YellowO - Yellow O
	{0, 255, 0, 255},     // GreenS - Green S
	{128, 0, 128, 255},   // PurpleT - Purple T
	{255, 0, 0, 255},     // RedZ - Red Z
	{128, 128, 128, 255}, // Locked - Gray
}

// Renderer handles the game's rendering
type Renderer struct {
	game     *tetris.Game
	boardImg *ebiten.Image
	font     font.Face
	logoImg  *ebiten.Image
}

// NewRenderer creates a new renderer for the game
func NewRenderer(game *tetris.Game) *Renderer {
	// Try to load the Tetris logo
	logoImg := loadImage("assets/tetris_logo.png")

	return &Renderer{
		game:     game,
		boardImg: ebiten.NewImage(ScreenWidth, ScreenHeight),
		font:     basicfont.Face7x13,
		logoImg:  logoImg,
	}
}

// Draw renders the game to the screen
func (r *Renderer) Draw(screen *ebiten.Image) {
	// Clear the screen
	screen.Fill(color.RGBA{
		R: (BackgroundColor >> 24) & 0xFF,
		G: (BackgroundColor >> 16) & 0xFF,
		B: (BackgroundColor >> 8) & 0xFF,
		A: BackgroundColor & 0xFF,
	})

	switch r.game.State {
	case tetris.StateMainMenu:
		r.drawMainMenu(screen)
	case tetris.StateMultiplayerSetup:
		r.drawMultiplayerSetup(screen)
	case tetris.StateMatchmaking:
		r.drawMatchmaking(screen)
	case tetris.StatePlaying, tetris.StatePaused:
		r.drawGame(screen)
		if r.game.State == tetris.StatePaused {
			r.drawPauseOverlay(screen)
		}
	case tetris.StateGameOver:
		r.drawGame(screen)
		r.drawGameOverOverlay(screen)
	}
}

// drawMainMenu draws the main menu with game mode options
func (r *Renderer) drawMainMenu(screen *ebiten.Image) {
	// Draw the Tetris logo if available
	if r.logoImg != nil {
		logoWidth := r.logoImg.Bounds().Dx()
		logoHeight := r.logoImg.Bounds().Dy()

		// Center the logo horizontally
		x := (ScreenWidth - logoWidth) / 2
		y := ScreenHeight/4 - logoHeight/2

		// Draw the logo
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(r.logoImg, op)
	} else {
		// Fallback to text if logo is not available
		msg := "TETRIS"
		x := (ScreenWidth - len(msg)*7) / 2
		y := ScreenHeight / 4
		text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility
	}

	msg := "Official Guidelines Edition"
	x := (ScreenWidth - len(msg)*7) / 2
	y := ScreenHeight / 3
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Menu options
	msg = "1. Single Player"
	x = (ScreenWidth - len(msg)*7) / 2
	y = ScreenHeight/2 - 10
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "2. Multiplayer"
	x = (ScreenWidth - len(msg)*7) / 2
	y = ScreenHeight/2 + 10
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "ESC. Quit"
	x = (ScreenWidth - len(msg)*7) / 2
	y = ScreenHeight/2 + 30
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "Controls:"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 40
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "Arrow Keys: Move"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "Up: Rotate"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "Space: Hard Drop"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "ESC: Pause"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "Shift: Hold Piece"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility
}

// drawGame draws the main game screen
func (r *Renderer) drawGame(screen *ebiten.Image) {
	// Draw the board frame
	vector.DrawFilledRect(
		screen,
		float32(BoardX-2),
		float32(BoardY-2),
		float32(tetris.BoardWidth*CellSize+4),
		float32(tetris.BoardHeight*CellSize+4),
		color.White,
		false,
	)

	// Draw the board cells - only draw the visible part (not the buffer rows)
	for y := 0; y < tetris.BoardHeight; y++ {
		// Calculate the actual y position in the board array (skip buffer rows)
		boardY := y + (tetris.BoardHeightWithBuffer - tetris.BoardHeight)

		for x := 0; x < tetris.BoardWidth; x++ {
			cell := r.game.Board.Cells[boardY][x]
			if cell != tetris.Empty {
				r.drawCell(screen, BoardX+x*CellSize, BoardY+y*CellSize, pieceColors[cell])
			} else {
				// Draw empty cell
				r.drawEmptyCell(screen, BoardX+x*CellSize, BoardY+y*CellSize)
			}
		}
	}

	// Draw the current piece
	if r.game.CurrentPiece != nil && r.game.State == tetris.StatePlaying {
		piece := r.game.CurrentPiece
		pieceColor := pieceColors[piece.Type]

		for i := 0; i < len(piece.Shape); i++ {
			for j := 0; j < len(piece.Shape[i]); j++ {
				if piece.Shape[i][j] {
					// Calculate the visual position, accounting for buffer rows
					visualY := piece.Y - (tetris.BoardHeightWithBuffer - tetris.BoardHeight)

					// Only draw if the piece is in the visible area
					if visualY+i >= 0 {
						x := BoardX + (piece.X+j)*CellSize
						y := BoardY + (visualY+i)*CellSize
						r.drawCell(screen, x, y, pieceColor)
					}
				}
			}
		}

		// Get ghost piece position
		ghostY := r.game.GetGhostPieceY()

		// Only draw ghost if it's different from the current position
		if ghostY > piece.Y {
			ghostColor := color.RGBA{pieceColor.R, pieceColor.G, pieceColor.B, 100}

			// Calculate the visual position for the ghost piece
			visualGhostY := ghostY - (tetris.BoardHeightWithBuffer - tetris.BoardHeight)

			for i := 0; i < len(piece.Shape); i++ {
				for j := 0; j < len(piece.Shape[i]); j++ {
					if piece.Shape[i][j] && visualGhostY+i >= 0 {
						x := BoardX + (piece.X+j)*CellSize
						y := BoardY + (visualGhostY+i)*CellSize
						r.drawCell(screen, x, y, ghostColor)
					}
				}
			}
		}
	}

	// Draw the next piece preview
	r.drawNextPiecePreview(screen)

	// Draw the held piece
	r.drawHeldPiece(screen)

	// Draw game stats
	r.drawGameStats(screen)
}

// drawCell draws a colored cell
func (r *Renderer) drawCell(screen *ebiten.Image, x, y int, clr color.RGBA) {
	// Draw the cell fill
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(CellSize),
		float32(CellSize),
		clr,
		false,
	)

	// Draw cell border (single call with a darker color)
	borderColor := color.RGBA{
		R: uint8(float64(clr.R) * 0.7),
		G: uint8(float64(clr.G) * 0.7),
		B: uint8(float64(clr.B) * 0.7),
		A: 255,
	}

	// Top border
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(CellSize),
		1,
		borderColor,
		false,
	)

	// Left border
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		1,
		float32(CellSize),
		borderColor,
		false,
	)

	// Right border
	vector.DrawFilledRect(
		screen,
		float32(x+CellSize-1),
		float32(y),
		1,
		float32(CellSize),
		borderColor,
		false,
	)

	// Bottom border
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y+CellSize-1),
		float32(CellSize),
		1,
		borderColor,
		false,
	)
}

// drawEmptyCell draws an empty cell
func (r *Renderer) drawEmptyCell(screen *ebiten.Image, x, y int) {
	// Draw the cell fill
	cellColor := color.RGBA{40, 40, 40, 255}
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(CellSize),
		float32(CellSize),
		cellColor,
		false,
	)

	// Draw cell border
	borderColor := color.RGBA{20, 20, 20, 255}

	// Top border
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(CellSize),
		1,
		borderColor,
		false,
	)

	// Left border
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		1,
		float32(CellSize),
		borderColor,
		false,
	)

	// Right border
	vector.DrawFilledRect(
		screen,
		float32(x+CellSize-1),
		float32(y),
		1,
		float32(CellSize),
		borderColor,
		false,
	)

	// Bottom border
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y+CellSize-1),
		float32(CellSize),
		1,
		borderColor,
		false,
	)
}

// drawNextPiecePreview draws the next piece preview
func (r *Renderer) drawNextPiecePreview(screen *ebiten.Image) {
	// Draw preview box with border
	vector.DrawFilledRect(
		screen,
		float32(PreviewX-12),
		float32(PreviewY-32),
		124,
		104,
		color.RGBA{100, 100, 100, 255},
		false,
	)

	vector.DrawFilledRect(
		screen,
		float32(PreviewX-10),
		float32(PreviewY-30),
		120,
		100,
		color.RGBA{60, 60, 60, 255},
		false,
	)

	text.Draw(screen, "Next Piece:", r.font, PreviewX-5, PreviewY-15, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	if r.game.NextPiece != nil {
		piece := r.game.NextPiece
		pieceColor := pieceColors[piece.Type]

		// Center the piece in the preview box
		offsetX := (4 - len(piece.Shape[0])) / 2 * CellSize
		offsetY := (4 - len(piece.Shape)) / 2 * CellSize

		for i := 0; i < len(piece.Shape); i++ {
			for j := 0; j < len(piece.Shape[i]); j++ {
				if piece.Shape[i][j] {
					x := PreviewX + j*CellSize + offsetX
					y := PreviewY + i*CellSize + offsetY
					r.drawCell(screen, x, y, pieceColor)
				}
			}
		}
	}
}

// drawHeldPiece draws the held piece
func (r *Renderer) drawHeldPiece(screen *ebiten.Image) {
	// Draw hold box with border
	vector.DrawFilledRect(
		screen,
		float32(HoldX-12),
		float32(HoldY-32),
		124,
		104,
		color.RGBA{100, 100, 100, 255},
		false,
	)

	vector.DrawFilledRect(
		screen,
		float32(HoldX-10),
		float32(HoldY-30),
		120,
		100,
		color.RGBA{60, 60, 60, 255},
		false,
	)

	text.Draw(screen, "Hold Piece:", r.font, HoldX-5, HoldY-15, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	heldPiece := r.game.GetHeldPiece()
	if heldPiece != nil {
		pieceColor := pieceColors[heldPiece.Type]

		// Center the piece in the hold box
		offsetX := (4 - len(heldPiece.Shape[0])) / 2 * CellSize
		offsetY := (4 - len(heldPiece.Shape)) / 2 * CellSize

		for i := 0; i < len(heldPiece.Shape); i++ {
			for j := 0; j < len(heldPiece.Shape[i]); j++ {
				if heldPiece.Shape[i][j] {
					x := HoldX + j*CellSize + offsetX
					y := HoldY + i*CellSize + offsetY
					r.drawCell(screen, x, y, pieceColor)
				}
			}
		}
	} else {
		// Draw empty indicator
		text.Draw(screen, "Empty", r.font, HoldX+5, HoldY+30, color.White) // nolint:staticcheck // Using deprecated API for compatibility
	}
}

// drawGameStats draws the game statistics
func (r *Renderer) drawGameStats(screen *ebiten.Image) {
	// Calculate box height based on multiplayer mode
	boxHeight := float32(160) // Default height
	if r.game.MultiplayerMode {
		boxHeight = 240 // Taller for player names
	}

	// Draw stats box with border
	vector.DrawFilledRect(
		screen,
		float32(PreviewX-12),
		float32(PreviewY+78),
		124,
		boxHeight+4, // Border height
		color.RGBA{100, 100, 100, 255},
		false,
	)

	vector.DrawFilledRect(
		screen,
		float32(PreviewX-10),
		float32(PreviewY+80),
		120,
		boxHeight, // Inner box height
		color.RGBA{60, 60, 60, 255},
		false,
	)

	// Show player names in multiplayer mode
	if r.game.MultiplayerMode {
		// Your name (top)
		yourName := r.game.UsernameInput
		if yourName == "" {
			yourName = "You"
		}
		text.Draw(screen, "You:", r.font, PreviewX-5, PreviewY+95, color.RGBA{255, 255, 0, 255}) // Yellow
		text.Draw(screen, yourName, r.font, PreviewX+5, PreviewY+110, color.RGBA{255, 255, 0, 255})

		// Opponent name
		opponentName := r.game.OpponentName
		if opponentName == "" {
			opponentName = "Opponent"
		}
		text.Draw(screen, "Opponent:", r.font, PreviewX-5, PreviewY+130, color.RGBA{255, 100, 100, 255}) // Light red
		text.Draw(screen, opponentName, r.font, PreviewX+5, PreviewY+145, color.RGBA{255, 100, 100, 255})
	}

	// Draw score
	text.Draw(screen, "Score:", r.font, PreviewX-5, PreviewY+170, color.White)                             // nolint:staticcheck // Using deprecated API for compatibility
	text.Draw(screen, fmt.Sprintf("%d", r.game.GetScore()), r.font, PreviewX+5, PreviewY+190, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Draw level
	text.Draw(screen, "Level:", r.font, PreviewX-5, PreviewY+210, color.White)                             // nolint:staticcheck // Using deprecated API for compatibility
	text.Draw(screen, fmt.Sprintf("%d", r.game.GetLevel()), r.font, PreviewX+5, PreviewY+230, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Draw lines cleared
	text.Draw(screen, "Lines:", r.font, PreviewX-5, PreviewY+250, color.White)                                    // nolint:staticcheck // Using deprecated API for compatibility
	text.Draw(screen, fmt.Sprintf("%d", r.game.GetLinesCleared()), r.font, PreviewX+5, PreviewY+270, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Draw Back-to-Back status
	if r.game.GetBackToBack() {
		text.Draw(screen, "Back-to-Back", r.font, PreviewX-5, PreviewY+290, color.RGBA{255, 215, 0, 255}) // Gold color
	}

	// Draw last clear type
	if r.game.GetLastClearWasTSpin() {
		text.Draw(screen, "T-Spin!", r.font, PreviewX-5, PreviewY+310, color.RGBA{255, 105, 180, 255}) // Hot pink
	}
}

// drawPauseOverlay draws the pause screen overlay
func (r *Renderer) drawPauseOverlay(screen *ebiten.Image) {
	// Semi-transparent overlay
	vector.DrawFilledRect(
		screen,
		0,
		0,
		float32(ScreenWidth),
		float32(ScreenHeight),
		color.RGBA{0, 0, 0, 128},
		false,
	)

	// Pause text
	msg := "PAUSED"
	x := (ScreenWidth - len(msg)*7) / 2
	y := ScreenHeight/2 - 10
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "Press ESC to resume"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 30
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility
}

// drawGameOverOverlay draws the game over screen overlay
func (r *Renderer) drawGameOverOverlay(screen *ebiten.Image) {
	// Semi-transparent overlay
	vector.DrawFilledRect(
		screen,
		0,
		0,
		float32(ScreenWidth),
		float32(ScreenHeight),
		color.RGBA{0, 0, 0, 192},
		false,
	)

	// Game over text
	msg := "GAME OVER"
	x := (ScreenWidth - len(msg)*7) / 2
	y := ScreenHeight/2 - 30
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Final score
	msg = fmt.Sprintf("Final Score: %d", r.game.GetScore())
	x = (ScreenWidth - len(msg)*7) / 2
	y += 30
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Restart instructions
	msg = "Press ENTER to play again"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 30
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility
}

// drawMultiplayerSetup draws the multiplayer setup screen
func (r *Renderer) drawMultiplayerSetup(screen *ebiten.Image) {
	msg := "MULTIPLAYER SETUP"
	x := (ScreenWidth - len(msg)*7) / 2
	y := ScreenHeight / 4
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	msg = "Enter your username:"
	x = (ScreenWidth - len(msg)*7) / 2
	y = ScreenHeight/2 - 40
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Show username input with cursor
	username := r.game.UsernameInput
	if username == "" {
		username = "_"
	} else {
		username += "_" // Show cursor
	}
	x = (ScreenWidth - len(username)*7) / 2
	y = ScreenHeight/2 - 10
	text.Draw(screen, username, r.font, x, y, color.RGBA{255, 255, 0, 255}) // Yellow text

	// Show character count
	charCount := fmt.Sprintf("(%d/12)", len(r.game.UsernameInput))
	x = (ScreenWidth - len(charCount)*7) / 2
	y = ScreenHeight/2 + 10
	text.Draw(screen, charCount, r.font, x, y, color.RGBA{128, 128, 128, 255}) // Gray text

	// Show connection status if any
	if r.game.ConnectionStatus != "" {
		status := r.game.ConnectionStatus
		x = (ScreenWidth - len(status)*7) / 2
		y = ScreenHeight/2 + 30
		statusColor := color.RGBA{255, 255, 255, 255} // White
		if status == "Connection error occurred" ||
			status == "Username must be at least 2 characters" ||
			status == "Username too long (max 12 characters)" {
			statusColor = color.RGBA{255, 0, 0, 255} // Red for errors
		}
		text.Draw(screen, status, r.font, x, y, statusColor) // nolint:staticcheck // Using deprecated API for compatibility
	}

	msg = "ENTER to connect | ESC to back"
	x = (ScreenWidth - len(msg)*7) / 2
	y = ScreenHeight/2 + 60
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility
}

// drawMatchmaking draws the matchmaking screen
func (r *Renderer) drawMatchmaking(screen *ebiten.Image) {
	msg := "FINDING MATCH..."
	x := (ScreenWidth - len(msg)*7) / 2
	y := ScreenHeight / 3
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Show connection status if available
	if r.game.ConnectionStatus != "" {
		status := r.game.ConnectionStatus
		x = (ScreenWidth - len(status)*7) / 2
		y = ScreenHeight / 2
		text.Draw(screen, status, r.font, x, y, color.RGBA{255, 255, 0, 255}) // Yellow text
	} else {
		msg = "Waiting for opponent..."
		x = (ScreenWidth - len(msg)*7) / 2
		y = ScreenHeight / 2
		text.Draw(screen, msg, r.font, x, y, color.RGBA{128, 128, 128, 255}) // Gray text
	}

	msg = "ESC to cancel"
	x = (ScreenWidth - len(msg)*7) / 2
	y = ScreenHeight/2 + 40
	text.Draw(screen, msg, r.font, x, y, color.White) // nolint:staticcheck // Using deprecated API for compatibility
}
