//go:build !headless
// +build !headless

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/briancain/go-tetris/internal/tetris"
	"github.com/briancain/go-tetris/internal/ui"
)

func main() {
	// Seed the random number generator (for Go 1.20+)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a new game
	game := tetris.NewGame()

	// Create a renderer
	renderer := ui.NewRenderer(game)

	// Create a new app
	app := tetris.NewApp(game, renderer)

	// Set up Ebiten
	ebiten.SetWindowSize(ui.ScreenWidth, ui.ScreenHeight)
	ebiten.SetWindowTitle("Go Tetris")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	// Run the game
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
