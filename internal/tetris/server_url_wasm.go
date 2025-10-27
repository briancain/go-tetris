//go:build js && wasm
// +build js,wasm

package tetris

import (
	"log"
	"syscall/js"
)

// getServerURL returns the server URL from JavaScript configuration for WebAssembly builds
func getServerURL() string {
	defer func() {
		if r := recover(); r != nil {
			// If JavaScript access fails, fall back to localhost
			log.Printf("Failed to read server URL from JavaScript, using localhost: %v", r)
		}
	}()

	config := js.Global().Get("TETRIS_CONFIG")
	if !config.IsUndefined() {
		serverURL := config.Get("serverURL")
		if !serverURL.IsUndefined() && serverURL.Type() == js.TypeString {
			url := serverURL.String()
			if url != "" && url != "{{SERVER_URL}}" {
				return url
			}
		}
	}

	// Default to localhost if JavaScript config is not available
	return "http://localhost:8080"
}
