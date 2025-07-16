// Package testutil provides utilities for testing
package testutil

import (
	"os"
)

func init() {
	// Set headless mode for Ebiten
	os.Setenv("EBITEN_HEADLESS", "1")
}
