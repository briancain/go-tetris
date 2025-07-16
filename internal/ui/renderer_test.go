package ui

import (
	"image/color"
	"testing"

	"github.com/briancain/go-tetris/internal/tetris"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewRenderer(t *testing.T) {
	game := tetris.NewGame()
	renderer := NewRenderer(game)

	if renderer == nil {
		t.Error("NewRenderer returned nil")
		return
	}

	if renderer.game != game {
		t.Error("Renderer's game reference is incorrect")
	}

	if renderer.boardImg == nil {
		t.Error("Renderer's boardImg is nil")
	}

	if renderer.font == nil {
		t.Error("Renderer's font is nil")
	}
}

func TestDrawCell(_ *testing.T) {
	game := tetris.NewGame()
	renderer := NewRenderer(game)

	// Create a test image
	img := ebiten.NewImage(100, 100)

	// Draw a cell
	testColor := color.RGBA{255, 0, 0, 255}
	renderer.drawCell(img, 10, 10, testColor)

	// We can't easily test the pixel colors in a unit test without complex mocking,
	// so we'll just verify that the method doesn't panic
}

func TestDrawEmptyCell(_ *testing.T) {
	game := tetris.NewGame()
	renderer := NewRenderer(game)

	// Create a test image
	img := ebiten.NewImage(100, 100)

	// Draw an empty cell
	renderer.drawEmptyCell(img, 10, 10)

	// Again, we're just verifying that the method doesn't panic
}

func TestScreenDimensions(t *testing.T) {
	// Test that screen dimensions are reasonable
	if ScreenWidth <= 0 || ScreenHeight <= 0 {
		t.Errorf("Invalid screen dimensions: %dx%d", ScreenWidth, ScreenHeight)
	}

	// Test that cell size is reasonable
	if CellSize <= 0 {
		t.Errorf("Invalid cell size: %d", CellSize)
	}

	// Test that board position is within screen bounds
	if BoardX < 0 || BoardX >= ScreenWidth || BoardY < 0 || BoardY >= ScreenHeight {
		t.Errorf("Board position out of bounds: (%d,%d)", BoardX, BoardY)
	}

	// Test that preview position is within screen bounds
	if PreviewX < 0 || PreviewX >= ScreenWidth || PreviewY < 0 || PreviewY >= ScreenHeight {
		t.Errorf("Preview position out of bounds: (%d,%d)", PreviewX, PreviewY)
	}

	// Test that hold position is within screen bounds
	if HoldX < 0 || HoldX >= ScreenWidth || HoldY < 0 || HoldY >= ScreenHeight {
		t.Errorf("Hold position out of bounds: (%d,%d)", HoldX, HoldY)
	}
}

func TestPieceColors(t *testing.T) {
	// Test that we have colors for all piece types
	if len(pieceColors) != 9 { // Empty + 7 piece types + Locked
		t.Errorf("Expected 9 piece colors, got %d", len(pieceColors))
	}

	// Test that colors are valid
	for i, c := range pieceColors {
		if i > 0 && c.A == 0 {
			t.Errorf("Piece color %d has zero alpha", i)
		}
	}
}
