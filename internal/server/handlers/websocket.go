package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/briancain/go-tetris/internal/server/services"
	"github.com/briancain/go-tetris/pkg/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	wsManager   *services.WebSocketManager
	authService *services.AuthService
	gameManager *services.GameManager
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(
	wsManager *services.WebSocketManager,
	authService *services.AuthService,
	gameManager *services.GameManager,
) *WebSocketHandler {
	return &WebSocketHandler{
		wsManager:   wsManager,
		authService: authService,
		gameManager: gameManager,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get session token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Validate token
	player, err := h.authService.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}

	// Add connection to manager
	h.wsManager.AddConnection(player.ID, conn)

	// Handle messages
	go h.handleMessages(player.ID, conn)
}

// handleMessages processes incoming WebSocket messages
func (h *WebSocketHandler) handleMessages(playerID string, conn *websocket.Conn) {
	defer func() {
		h.handlePlayerDisconnect(playerID)
		h.wsManager.RemoveConnection(playerID)
		conn.Close()
	}()

	for {
		_, messageData, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error for player %s: %v", playerID, err)
			break
		}

		// Parse message
		var message map[string]interface{}
		err = json.Unmarshal(messageData, &message)
		if err != nil {
			log.Printf("Failed to parse WebSocket message: %v", err)
			continue
		}

		// Handle different message types
		messageType, ok := message["type"].(string)
		if !ok {
			log.Printf("Missing message type")
			continue
		}

		switch messageType {
		case "game_move":
			h.handleGameMove(playerID, message)
		case "game_state":
			h.handleGameState(playerID, message)
		case "game_over":
			h.handleGameOver(playerID, message)
		case "rematch_request":
			h.handleRematchRequest(playerID, message)
		case "ping":
			h.handlePing(playerID)
		default:
			log.Printf("Unknown message type: %s", messageType)
		}
	}
}

// handleGameMove processes a game move message
func (h *WebSocketHandler) handleGameMove(playerID string, message map[string]interface{}) {
	moveType, ok := message["moveType"].(string)
	if !ok {
		return
	}

	move := &models.GameMove{
		PlayerID:  playerID,
		MoveType:  moveType,
		Timestamp: time.Now(),
	}

	err := h.gameManager.HandleGameMove(playerID, move)
	if err != nil {
		log.Printf("Failed to handle game move: %v", err)
	}
}

// handleGameState processes a game state update
func (h *WebSocketHandler) handleGameState(playerID string, message map[string]interface{}) {
	// Extract game state data
	boardInterface, _ := message["board"].([]interface{})
	score, _ := message["score"].(float64)
	level, _ := message["level"].(float64)
	lines, _ := message["lines"].(float64)

	// Convert board to [][]int
	var gameBoard [][]int
	for _, rowInterface := range boardInterface {
		row, ok := rowInterface.([]interface{})
		if !ok {
			continue
		}
		var gameRow []int
		for _, cell := range row {
			cellValue, ok := cell.(float64)
			if ok {
				gameRow = append(gameRow, int(cellValue))
			}
		}
		gameBoard = append(gameBoard, gameRow)
	}

	state := &models.GameState{
		PlayerID:  playerID,
		Board:     gameBoard,
		Score:     int(score),
		Level:     int(level),
		Lines:     int(lines),
		Timestamp: time.Now(),
	}

	err := h.gameManager.HandleGameState(playerID, state)
	if err != nil {
		log.Printf("Failed to handle game state: %v", err)
	}
}

// handleGameOver processes a game over message
func (h *WebSocketHandler) handleGameOver(playerID string, message map[string]interface{}) {
	gameID, ok := message["gameId"].(string)
	if !ok {
		return
	}

	err := h.gameManager.EndGame(gameID, playerID)
	if err != nil {
		log.Printf("Failed to handle game over: %v", err)
	}
}

// handleRematchRequest processes a rematch request
func (h *WebSocketHandler) handleRematchRequest(playerID string, message map[string]interface{}) {
	err := h.gameManager.HandleRematchRequest(playerID)
	if err != nil {
		log.Printf("Failed to handle rematch request: %v", err)
	}
}

// handlePing responds to ping messages
func (h *WebSocketHandler) handlePing(playerID string) {
	pongMsg := map[string]interface{}{
		"type": "pong",
	}

	data, _ := json.Marshal(pongMsg)
	h.wsManager.SendToPlayer(playerID, data)
}

// handlePlayerDisconnect handles when a player disconnects
func (h *WebSocketHandler) handlePlayerDisconnect(playerID string) {
	log.Printf("Player %s disconnected", playerID)

	// Notify game manager of disconnect
	err := h.gameManager.HandlePlayerDisconnect(playerID)
	if err != nil {
		log.Printf("Failed to handle player disconnect: %v", err)
	}
}
