package services

import (
	"sync"

	"github.com/gorilla/websocket"

	"github.com/briancain/go-tetris/internal/server/logger"
)

// connWrapper wraps a WebSocket connection with a mutex for safe concurrent writes
type connWrapper struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

// WebSocketManager handles WebSocket connections
type WebSocketManager struct {
	connections map[string]*connWrapper // playerID -> connection wrapper
	mu          sync.RWMutex
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[string]*connWrapper),
	}
}

// AddConnection adds a WebSocket connection for a player
func (wsm *WebSocketManager) AddConnection(playerID string, conn *websocket.Conn) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	// Close existing connection if any
	if existingWrapper, exists := wsm.connections[playerID]; exists {
		existingWrapper.conn.Close()
	}

	wsm.connections[playerID] = &connWrapper{conn: conn}
	logger.Logger.Info("WebSocket connection added",
		"playerID", playerID,
	)
}

// RemoveConnection removes a WebSocket connection
func (wsm *WebSocketManager) RemoveConnection(playerID string) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	if wrapper, exists := wsm.connections[playerID]; exists {
		wrapper.conn.Close()
		delete(wsm.connections, playerID)
		logger.Logger.Info("WebSocket connection removed",
			"playerID", playerID,
		)
	}
}

// SendToPlayer sends a message to a specific player
func (wsm *WebSocketManager) SendToPlayer(playerID string, message []byte) {
	wsm.mu.RLock()
	wrapper, exists := wsm.connections[playerID]
	wsm.mu.RUnlock()

	if !exists {
		logger.Logger.Warn("No WebSocket connection found for player",
			"playerID", playerID,
		)
		return
	}

	// Use per-connection mutex to prevent concurrent writes
	wrapper.mu.Lock()
	err := wrapper.conn.WriteMessage(websocket.TextMessage, message)
	wrapper.mu.Unlock()

	if err != nil {
		logger.Logger.Error("Failed to send WebSocket message",
			"playerID", playerID,
			"error", err,
		)
		wsm.RemoveConnection(playerID)
	}
}

// BroadcastToAll sends a message to all connected players
func (wsm *WebSocketManager) BroadcastToAll(message []byte) {
	wsm.mu.RLock()
	defer wsm.mu.RUnlock()

	for playerID, wrapper := range wsm.connections {
		wrapper.mu.Lock()
		err := wrapper.conn.WriteMessage(websocket.TextMessage, message)
		wrapper.mu.Unlock()
		if err != nil {
			logger.Logger.Error("Failed to broadcast WebSocket message",
				"playerID", playerID,
				"error", err,
			)
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

// Shutdown gracefully closes all WebSocket connections
func (wsm *WebSocketManager) Shutdown() {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	logger.Logger.Info("Shutting down WebSocket connections", "count", len(wsm.connections))

	for playerID, wrapper := range wsm.connections {
		// Use per-connection mutex for safe close message
		wrapper.mu.Lock()
		// Send close message to client
		err := wrapper.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server shutting down"))
		if err != nil {
			logger.Logger.Warn("Failed to send close message", "playerID", playerID, "error", err)
		}

		// Close the connection
		err = wrapper.conn.Close()
		wrapper.mu.Unlock()
		if err != nil {
			logger.Logger.Warn("Failed to close WebSocket connection", "playerID", playerID, "error", err)
		}
	}

	// Clear all connections
	wsm.connections = make(map[string]*connWrapper)
	logger.Logger.Info("All WebSocket connections closed")
}
