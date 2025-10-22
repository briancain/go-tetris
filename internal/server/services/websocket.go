package services

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketManager handles WebSocket connections
type WebSocketManager struct {
	connections map[string]*websocket.Conn // playerID -> connection
	mu          sync.RWMutex
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[string]*websocket.Conn),
	}
}

// AddConnection adds a WebSocket connection for a player
func (wsm *WebSocketManager) AddConnection(playerID string, conn *websocket.Conn) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	// Close existing connection if any
	if existingConn, exists := wsm.connections[playerID]; exists {
		existingConn.Close()
	}

	wsm.connections[playerID] = conn
	log.Printf("Added WebSocket connection for player %s", playerID)
}

// RemoveConnection removes a WebSocket connection
func (wsm *WebSocketManager) RemoveConnection(playerID string) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	if conn, exists := wsm.connections[playerID]; exists {
		conn.Close()
		delete(wsm.connections, playerID)
		log.Printf("Removed WebSocket connection for player %s", playerID)
	}
}

// SendToPlayer sends a message to a specific player
func (wsm *WebSocketManager) SendToPlayer(playerID string, message []byte) {
	wsm.mu.RLock()
	conn, exists := wsm.connections[playerID]
	wsm.mu.RUnlock()

	if !exists {
		log.Printf("No WebSocket connection for player %s", playerID)
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Failed to send message to player %s: %v", playerID, err)
		wsm.RemoveConnection(playerID)
	}
}

// BroadcastToAll sends a message to all connected players
func (wsm *WebSocketManager) BroadcastToAll(message []byte) {
	wsm.mu.RLock()
	defer wsm.mu.RUnlock()

	for playerID, conn := range wsm.connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Failed to broadcast to player %s: %v", playerID, err)
			go wsm.RemoveConnection(playerID)
		}
	}
}

// GetConnectionCount returns the number of active connections
func (wsm *WebSocketManager) GetConnectionCount() int {
	wsm.mu.RLock()
	defer wsm.mu.RUnlock()
	return len(wsm.connections)
}
