//go:build headless
// +build headless

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/briancain/go-tetris/internal/tetris"
)

func main() {
	// Seed the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// For CI builds, we use a headless game implementation
	// that doesn't require a display
	headlessGame := &tetris.HeadlessGame{}

	// This is just to make the build succeed in CI
	// The actual game won't run in headless mode
	if err := ebiten.RunGame(headlessGame); err != nil {
		log.Fatal(err)
	}
}
