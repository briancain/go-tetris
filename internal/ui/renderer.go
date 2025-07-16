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
	case tetris.StateMenu:
		r.drawMenu(screen)
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

// drawMenu draws the main menu
func (r *Renderer) drawMenu(screen *ebiten.Image) {
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

	msg = "Press ENTER to start"
	x = (ScreenWidth - len(msg)*7) / 2
	y = ScreenHeight/2 + 20
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

		// Draw ghost piece (preview of where the piece will land)
		ghostY := piece.Y
		for {
			testY := ghostY + 1
			testPiece := piece.Copy()
			testPiece.Y = testY

			if !r.game.Board.IsValidPosition(testPiece, testPiece.X, testY) {
				break
			}

			ghostY = testY
		}

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
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(CellSize),
		float32(CellSize),
		clr,
		false,
	)

	// Draw cell border
	borderColor := color.RGBA{0, 0, 0, 255}
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(CellSize),
		1,
		borderColor,
		false,
	)
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		1,
		float32(CellSize),
		borderColor,
		false,
	)
	vector.DrawFilledRect(
		screen,
		float32(x+CellSize-1),
		float32(y),
		1,
		float32(CellSize),
		borderColor,
		false,
	)
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
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(CellSize),
		float32(CellSize),
		color.RGBA{40, 40, 40, 255},
		false,
	)

	// Draw cell border
	borderColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(CellSize),
		1,
		borderColor,
		false,
	)
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		1,
		float32(CellSize),
		borderColor,
		false,
	)
	vector.DrawFilledRect(
		screen,
		float32(x+CellSize-1),
		float32(y),
		1,
		float32(CellSize),
		borderColor,
		false,
	)
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
	// Draw stats box with border
	vector.DrawFilledRect(
		screen,
		float32(PreviewX-12),
		float32(PreviewY+78),
		124,
		164, // Increased height to accommodate more stats
		color.RGBA{100, 100, 100, 255},
		false,
	)

	vector.DrawFilledRect(
		screen,
		float32(PreviewX-10),
		float32(PreviewY+80),
		120,
		160, // Increased height to accommodate more stats
		color.RGBA{60, 60, 60, 255},
		false,
	)

	// Draw score
	text.Draw(screen, "Score:", r.font, PreviewX-5, PreviewY+100, color.White)                             // nolint:staticcheck // Using deprecated API for compatibility
	text.Draw(screen, fmt.Sprintf("%d", r.game.GetScore()), r.font, PreviewX+5, PreviewY+120, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Draw level
	text.Draw(screen, "Level:", r.font, PreviewX-5, PreviewY+140, color.White)                             // nolint:staticcheck // Using deprecated API for compatibility
	text.Draw(screen, fmt.Sprintf("%d", r.game.GetLevel()), r.font, PreviewX+5, PreviewY+160, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Draw lines cleared
	text.Draw(screen, "Lines:", r.font, PreviewX-5, PreviewY+180, color.White)                                    // nolint:staticcheck // Using deprecated API for compatibility
	text.Draw(screen, fmt.Sprintf("%d", r.game.GetLinesCleared()), r.font, PreviewX+5, PreviewY+200, color.White) // nolint:staticcheck // Using deprecated API for compatibility

	// Draw Back-to-Back status
	if r.game.GetBackToBack() {
		text.Draw(screen, "Back-to-Back", r.font, PreviewX-5, PreviewY+220, color.RGBA{255, 215, 0, 255}) // Gold color
	}

	// Draw last clear type
	if r.game.GetLastClearWasTSpin() {
		text.Draw(screen, "T-Spin!", r.font, PreviewX-5, PreviewY+240, color.RGBA{255, 105, 180, 255}) // Hot pink
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
