package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketManager_Shutdown(t *testing.T) {
	wsManager := NewWebSocketManager()

	// Create mock WebSocket connections
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		// Keep connection alive
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}))
	defer server.Close()

	// Connect multiple clients
	clients := make([]*websocket.Conn, 3)
	for i := 0; i < 3; i++ {
		url := "ws" + strings.TrimPrefix(server.URL, "http")
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("Failed to connect client %d: %v", i, err)
		}
		clients[i] = conn

		// Add to manager
		playerID := "player" + string(rune('1'+i))
		wsManager.AddConnection(playerID, conn)
	}

	// Verify connections are added
	if count := wsManager.GetConnectionCount(); count != 3 {
		t.Errorf("Expected 3 connections, got %d", count)
	}

	// Shutdown
	wsManager.Shutdown()

	// Verify all connections are closed
	if count := wsManager.GetConnectionCount(); count != 0 {
		t.Errorf("Expected 0 connections after shutdown, got %d", count)
	}

	// Verify clients detect connection closure (any close error is acceptable)
	for i, client := range clients {
		client.SetReadDeadline(time.Now().Add(1 * time.Second))
		_, _, err := client.ReadMessage()
		if err == nil {
			t.Errorf("Client %d should have detected connection closure", i)
		}
		// Any connection error indicates the connection was closed
	}
}

func TestWebSocketManager_ShutdownEmptyConnections(t *testing.T) {
	wsManager := NewWebSocketManager()

	// Shutdown with no connections should not panic
	wsManager.Shutdown()

	// Verify count is still 0
	if count := wsManager.GetConnectionCount(); count != 0 {
		t.Errorf("Expected 0 connections, got %d", count)
	}
}

func TestWebSocketManager_ShutdownWithFailedClose(t *testing.T) {
	wsManager := NewWebSocketManager()

	// Create a mock connection that will fail to close
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		// Immediately close server side to cause client close to fail
		conn.Close()
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	wsManager.AddConnection("player1", conn)

	// This should not panic even if close fails
	wsManager.Shutdown()

	// Verify connection is removed
	if count := wsManager.GetConnectionCount(); count != 0 {
		t.Errorf("Expected 0 connections after shutdown, got %d", count)
	}
}
