//go:build headless
// +build headless

package tetris

// This file provides a minimal implementation for CI and cross-compilation builds
// It allows the game to be built without requiring a display or CGO

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// HeadlessDriver provides a minimal implementation for CI builds
// that doesn't require a display or graphics libraries
type HeadlessDriver struct{}

// IsScreenClearedEveryFrame returns true to indicate the screen is cleared every frame
func (d *HeadlessDriver) IsScreenClearedEveryFrame() bool {
	return true
}

// SetScreenClearedEveryFrame is a no-op in headless mode
func (d *HeadlessDriver) SetScreenClearedEveryFrame(cleared bool) {
	// No-op
}

// IsVsyncEnabled returns false in headless mode
func (d *HeadlessDriver) IsVsyncEnabled() bool {
	return false
}

// SetVsyncEnabled is a no-op in headless mode
func (d *HeadlessDriver) SetVsyncEnabled(enabled bool) {
	// No-op
}

// DeviceScaleFactor returns 1.0 in headless mode
func (d *HeadlessDriver) DeviceScaleFactor() float64 {
	return 1.0
}

// Run is a no-op in headless mode
func (d *HeadlessDriver) Run(game ebiten.Game) error {
	// In headless mode, we don't actually run the game loop
	return nil
}
