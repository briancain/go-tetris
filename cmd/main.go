package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/briancain/go-tetris/internal/tetris"
	"github.com/briancain/go-tetris/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Seed the random number generator (for Go 1.20+)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create the game
	game := tetris.NewGame()
	
	// Create the renderer
	renderer := ui.NewRenderer(game)
	
	// Create the game application
	tetrisGame := tetris.NewTetrisGame(game, renderer)
	
	// Set up the window
	ebiten.SetWindowSize(ui.ScreenWidth, ui.ScreenHeight)
	ebiten.SetWindowTitle("Go Tetris")
	
	// Run the game
	if err := ebiten.RunGame(tetrisGame); err != nil {
		log.Fatal(err)
	}
}
