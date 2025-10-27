//go:build !js || !wasm
// +build !js !wasm

package tetris

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Connect establishes WebSocket connection
func (mc *MultiplayerClient) Connect() error {
	if mc.sessionToken == "" {
		return fmt.Errorf("must login before connecting")
	}

	u, err := url.Parse(mc.serverURL)
	if err != nil {
		return err
	}

	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	u.Path = "/ws"
	u.RawQuery = "token=" + mc.sessionToken

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	mc.conn = conn
	mc.connected = true

	// Start message reader
	go mc.readMessages()

	log.Printf("Multiplayer: Connected to server")
	return nil
}
