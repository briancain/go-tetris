//go:build headless
// +build headless

// This file is used for cross-compilation in CI environments
// It provides a minimal implementation that can be built without
// requiring a display or graphics libraries
// Use with build tags: headless,ebitennogl,ebitennonscreen

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/briancain/go-tetris/internal/tetris"
)

func main() {
	// Seed the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// For CI builds, we use a headless implementation
	// that doesn't require a display
	headlessGame := &tetris.HeadlessGame{}
	headlessDriver := &tetris.HeadlessDriver{}

	// In headless mode, we don't actually run the game
	// This is just to make the build succeed in CI
	if err := headlessDriver.Run(headlessGame); err != nil {
		log.Fatal(err)
	}

	log.Println("Headless build completed successfully")
}
