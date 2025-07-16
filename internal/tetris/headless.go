//go:build headless
// +build headless

package tetris

// This file contains headless implementations for CI builds
// It's only used during cross-compilation in CI environments

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// HeadlessGame is a minimal implementation for CI builds
type HeadlessGame struct{}

func (g *HeadlessGame) Update() error {
	return nil
}

func (g *HeadlessGame) Draw(screen *ebiten.Image) {
	// No-op for headless builds
}

func (g *HeadlessGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

// InitHeadlessMode can be called in main when the ebitenginedummy tag is set
func InitHeadlessMode() {
	// This function would be called in your main.go when the ebitenginedummy tag is set
	// It allows the binary to be built without requiring a display
}
