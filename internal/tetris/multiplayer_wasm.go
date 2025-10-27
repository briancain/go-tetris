//go:build js && wasm
// +build js,wasm

package tetris

import (
	"encoding/json"
	"fmt"
	"log"
	"syscall/js"
)

// Connect establishes WebSocket connection using browser WebSocket API
func (mc *MultiplayerClient) Connect() error {
	if mc.sessionToken == "" {
		return fmt.Errorf("must login before connecting")
	}

	// Get WebSocket URL from JavaScript config
	config := js.Global().Get("TETRIS_CONFIG")
	wsBaseURL := config.Get("wsURL").String()
	if wsBaseURL == "" || wsBaseURL == "{{WS_URL}}" {
		wsBaseURL = "ws://localhost:8080"
	}

	wsURL := wsBaseURL + "/ws?token=" + mc.sessionToken

	ws := js.Global().Get("WebSocket").New(wsURL)

	// Wait for connection
	openChan := make(chan error, 1)

	onOpen := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		openChan <- nil
		return nil
	})

	onError := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		openChan <- fmt.Errorf("websocket connection failed")
		return nil
	})

	ws.Set("onopen", onOpen)
	ws.Set("onerror", onError)

	err := <-openChan
	onOpen.Release()
	onError.Release()

	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	// Setup message handler
	onMessage := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			data := args[0].Get("data").String()
			var msg map[string]interface{}
			if err := json.Unmarshal([]byte(data), &msg); err == nil {
				select {
				case mc.messages <- msg:
				default:
				}
			}
		}
		return nil
	})
	ws.Set("onmessage", onMessage)

	mc.connected = true
	log.Printf("Multiplayer: Connected to server")
	return nil
}
