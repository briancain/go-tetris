//go:build ci
// +build ci

package ui

import (
	"os"
)

func init() {
	// Set environment variables to help with CI testing
	os.Setenv("EBITEN_HEADLESS", "1")
}
