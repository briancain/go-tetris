//go:build !js || !wasm
// +build !js !wasm

package tetris

// getServerURL returns localhost for non-WebAssembly builds
func getServerURL() string {
	return "http://localhost:8080"
}
