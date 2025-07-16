package ui

import (
	"fmt"
	"image/color"

	"github.com/briancain/go-tetris/internal/tetris"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
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
	{0, 255, 255, 255},   // I - Cyan
	{0, 0, 255, 255},     // J - Blue
	{255, 165, 0, 255},   // L - Orange
	{255, 255, 0, 255},   // O - Yellow
	{0, 255, 0, 255},     // S - Green
	{128, 0, 128, 255},   // T - Purple
	{255, 0, 0, 255},     // Z - Red
	{128, 128, 128, 255}, // Locked - Gray
}

// Renderer handles the game's rendering
type Renderer struct {
	game     *tetris.Game
	boardImg *ebiten.Image
	font     font.Face
}

// NewRenderer creates a new renderer for the game
func NewRenderer(game *tetris.Game) *Renderer {
	return &Renderer{
		game:     game,
		boardImg: ebiten.NewImage(ScreenWidth, ScreenHeight),
		font:     basicfont.Face7x13,
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
	msg := "TETRIS"
	x := (ScreenWidth - len(msg)*7) / 2
	y := ScreenHeight/3 - 20
	text.Draw(screen, msg, r.font, x, y, color.White)

	msg = "Press ENTER to start"
	x = (ScreenWidth - len(msg)*7) / 2
	y = ScreenHeight/2 + 20
	text.Draw(screen, msg, r.font, x, y, color.White)

	msg = "Controls:"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 40
	text.Draw(screen, msg, r.font, x, y, color.White)

	msg = "Arrow Keys: Move"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White)

	msg = "Up: Rotate"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White)

	msg = "Space: Hard Drop"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White)

	msg = "ESC: Pause"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White)

	msg = "Shift: Hold Piece"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 20
	text.Draw(screen, msg, r.font, x, y, color.White)
}

// drawGame draws the main game screen
func (r *Renderer) drawGame(screen *ebiten.Image) {
	// Draw the board frame
	ebitenutil.DrawRect(
		screen,
		float64(BoardX-2),
		float64(BoardY-2),
		float64(tetris.BoardWidth*CellSize+4),
		float64(tetris.BoardHeight*CellSize+4),
		color.White,
	)

	// Draw the board cells
	for y := 0; y < tetris.BoardHeight; y++ {
		for x := 0; x < tetris.BoardWidth; x++ {
			cell := r.game.Board.Cells[y][x]
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
					x := BoardX + (piece.X+j)*CellSize
					y := BoardY + (piece.Y+i)*CellSize
					r.drawCell(screen, x, y, pieceColor)
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
			for i := 0; i < len(piece.Shape); i++ {
				for j := 0; j < len(piece.Shape[i]); j++ {
					if piece.Shape[i][j] {
						x := BoardX + (piece.X+j)*CellSize
						y := BoardY + (ghostY+i)*CellSize
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
	ebitenutil.DrawRect(
		screen,
		float64(x),
		float64(y),
		float64(CellSize),
		float64(CellSize),
		clr,
	)

	// Draw cell border
	borderColor := color.RGBA{0, 0, 0, 255}
	ebitenutil.DrawRect(
		screen,
		float64(x),
		float64(y),
		float64(CellSize),
		1,
		borderColor,
	)
	ebitenutil.DrawRect(
		screen,
		float64(x),
		float64(y),
		1,
		float64(CellSize),
		borderColor,
	)
	ebitenutil.DrawRect(
		screen,
		float64(x+CellSize-1),
		float64(y),
		1,
		float64(CellSize),
		borderColor,
	)
	ebitenutil.DrawRect(
		screen,
		float64(x),
		float64(y+CellSize-1),
		float64(CellSize),
		1,
		borderColor,
	)
}

// drawEmptyCell draws an empty cell
func (r *Renderer) drawEmptyCell(screen *ebiten.Image, x, y int) {
	ebitenutil.DrawRect(
		screen,
		float64(x),
		float64(y),
		float64(CellSize),
		float64(CellSize),
		color.RGBA{40, 40, 40, 255},
	)

	// Draw cell border
	borderColor := color.RGBA{20, 20, 20, 255}
	ebitenutil.DrawRect(
		screen,
		float64(x),
		float64(y),
		float64(CellSize),
		1,
		borderColor,
	)
	ebitenutil.DrawRect(
		screen,
		float64(x),
		float64(y),
		1,
		float64(CellSize),
		borderColor,
	)
	ebitenutil.DrawRect(
		screen,
		float64(x+CellSize-1),
		float64(y),
		1,
		float64(CellSize),
		borderColor,
	)
	ebitenutil.DrawRect(
		screen,
		float64(x),
		float64(y+CellSize-1),
		float64(CellSize),
		1,
		borderColor,
	)
}

// drawNextPiecePreview draws the next piece preview
func (r *Renderer) drawNextPiecePreview(screen *ebiten.Image) {
	// Draw preview box with border
	ebitenutil.DrawRect(
		screen,
		float64(PreviewX-12),
		float64(PreviewY-32),
		124,
		104,
		color.RGBA{100, 100, 100, 255},
	)

	ebitenutil.DrawRect(
		screen,
		float64(PreviewX-10),
		float64(PreviewY-30),
		120,
		100,
		color.RGBA{60, 60, 60, 255},
	)

	text.Draw(screen, "Next Piece:", r.font, PreviewX-5, PreviewY-15, color.White)

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
	ebitenutil.DrawRect(
		screen,
		float64(HoldX-12),
		float64(HoldY-32),
		124,
		104,
		color.RGBA{100, 100, 100, 255},
	)

	ebitenutil.DrawRect(
		screen,
		float64(HoldX-10),
		float64(HoldY-30),
		120,
		100,
		color.RGBA{60, 60, 60, 255},
	)

	text.Draw(screen, "Hold Piece:", r.font, HoldX-5, HoldY-15, color.White)

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
		text.Draw(screen, "Empty", r.font, HoldX+5, HoldY+30, color.White)
	}
}

// drawGameStats draws the game statistics
func (r *Renderer) drawGameStats(screen *ebiten.Image) {
	// Draw stats box with border
	ebitenutil.DrawRect(
		screen,
		float64(PreviewX-12),
		float64(PreviewY+78),
		124,
		124,
		color.RGBA{100, 100, 100, 255},
	)

	ebitenutil.DrawRect(
		screen,
		float64(PreviewX-10),
		float64(PreviewY+80),
		120,
		120,
		color.RGBA{60, 60, 60, 255},
	)

	// Draw score
	text.Draw(screen, "Score:", r.font, PreviewX-5, PreviewY+100, color.White)
	text.Draw(screen, fmt.Sprintf("%d", r.game.GetScore()), r.font, PreviewX+5, PreviewY+120, color.White)

	// Draw level
	text.Draw(screen, "Level:", r.font, PreviewX-5, PreviewY+140, color.White)
	text.Draw(screen, fmt.Sprintf("%d", r.game.GetLevel()), r.font, PreviewX+5, PreviewY+160, color.White)

	// Draw lines cleared
	text.Draw(screen, "Lines:", r.font, PreviewX-5, PreviewY+180, color.White)
	text.Draw(screen, fmt.Sprintf("%d", r.game.GetLinesCleared()), r.font, PreviewX+5, PreviewY+200, color.White)
}

// drawPauseOverlay draws the pause screen overlay
func (r *Renderer) drawPauseOverlay(screen *ebiten.Image) {
	// Semi-transparent overlay
	ebitenutil.DrawRect(
		screen,
		0,
		0,
		float64(ScreenWidth),
		float64(ScreenHeight),
		color.RGBA{0, 0, 0, 128},
	)

	// Pause text
	msg := "PAUSED"
	x := (ScreenWidth - len(msg)*7) / 2
	y := ScreenHeight/2 - 10
	text.Draw(screen, msg, r.font, x, y, color.White)

	msg = "Press ESC to resume"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 30
	text.Draw(screen, msg, r.font, x, y, color.White)
}

// drawGameOverOverlay draws the game over screen overlay
func (r *Renderer) drawGameOverOverlay(screen *ebiten.Image) {
	// Semi-transparent overlay
	ebitenutil.DrawRect(
		screen,
		0,
		0,
		float64(ScreenWidth),
		float64(ScreenHeight),
		color.RGBA{0, 0, 0, 192},
	)

	// Game over text
	msg := "GAME OVER"
	x := (ScreenWidth - len(msg)*7) / 2
	y := ScreenHeight/2 - 30
	text.Draw(screen, msg, r.font, x, y, color.White)

	// Final score
	msg = fmt.Sprintf("Final Score: %d", r.game.GetScore())
	x = (ScreenWidth - len(msg)*7) / 2
	y += 30
	text.Draw(screen, msg, r.font, x, y, color.White)

	// Restart instructions
	msg = "Press ENTER to play again"
	x = (ScreenWidth - len(msg)*7) / 2
	y += 30
	text.Draw(screen, msg, r.font, x, y, color.White)
}
