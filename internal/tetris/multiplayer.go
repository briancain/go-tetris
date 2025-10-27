package tetris

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

// MultiplayerClient handles server communication
type MultiplayerClient struct {
	conn         *websocket.Conn
	serverURL    string
	sessionToken string
	playerID     string
	username     string
	gameID       string
	connected    bool
	messages     chan map[string]interface{}
}

// NewMultiplayerClient creates a new multiplayer client
func NewMultiplayerClient(serverURL string) *MultiplayerClient {
	return &MultiplayerClient{
		serverURL: serverURL,
		messages:  make(chan map[string]interface{}, 100),
	}
}

// Login authenticates with the server
func (mc *MultiplayerClient) Login(username string) error {
	// Prepare login request
	loginReq := map[string]string{
		"username": username,
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %v", err)
	}

	// Make HTTP request to login endpoint
	resp, err := http.Post(mc.serverURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s", strings.TrimSpace(string(body)))
	}

	// Parse response
	var loginResp struct {
		PlayerID     string `json:"playerId"`
		Username     string `json:"username"`
		SessionToken string `json:"sessionToken"`
	}

	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return fmt.Errorf("failed to parse login response: %v", err)
	}

	// Store authentication info
	mc.username = loginResp.Username
	mc.sessionToken = loginResp.SessionToken
	mc.playerID = loginResp.PlayerID

	log.Printf("Multiplayer: Logged in as %s (ID: %s)", mc.username, mc.playerID)
	return nil
}

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

// JoinQueue joins the matchmaking queue
func (mc *MultiplayerClient) JoinQueue() error {
	if mc.sessionToken == "" {
		return fmt.Errorf("not logged in")
	}

	// Make HTTP request to join queue
	req, err := http.NewRequest("POST", mc.serverURL+"/api/matchmaking/queue", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+mc.sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to join queue: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("join queue failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Multiplayer: Joined matchmaking queue")
	return nil
}

// SendGameMove sends a move to the server
func (mc *MultiplayerClient) SendGameMove(moveType string) error {
	if !mc.connected {
		return nil // Silently ignore if not connected
	}

	message := map[string]interface{}{
		"type":     "game_move",
		"moveType": moveType,
	}

	return mc.sendMessage(message)
}

// SendGameState sends current game state to server
func (mc *MultiplayerClient) SendGameState(board [][]Cell, score, level, lines int) error {
	if !mc.connected {
		return nil // Silently ignore if not connected
	}

	// Convert board to int array for JSON
	boardInt := make([][]int, len(board))
	for i, row := range board {
		boardInt[i] = make([]int, len(row))
		for j, cell := range row {
			boardInt[i][j] = int(cell)
		}
	}

	message := map[string]interface{}{
		"type":  "game_state",
		"board": boardInt,
		"score": score,
		"level": level,
		"lines": lines,
	}

	return mc.sendMessage(message)
}

// GetMessage returns the next message from the server (non-blocking)
func (mc *MultiplayerClient) GetMessage() map[string]interface{} {
	select {
	case msg := <-mc.messages:
		return msg
	default:
		return nil
	}
}

// Close closes the connection
func (mc *MultiplayerClient) Close() {
	if mc.conn != nil {
		mc.conn.Close()
		mc.connected = false
	}
}

// IsConnected returns whether the client is connected
func (mc *MultiplayerClient) IsConnected() bool {
	return mc.connected
}

// GetGameID returns the current game ID
func (mc *MultiplayerClient) GetGameID() string {
	return mc.gameID
}

// SendGameOver sends a game over message to the server
func (mc *MultiplayerClient) SendGameOver() error {
	if !mc.connected {
		return nil // Silently ignore if not connected
	}

	message := map[string]interface{}{
		"type":   "game_over",
		"gameId": mc.gameID,
	}

	return mc.sendMessage(message)
}

// GetUsername returns the username
func (mc *MultiplayerClient) GetUsername() string {
	return mc.username
}

// sendMessage sends a message via WebSocket
func (mc *MultiplayerClient) sendMessage(message map[string]interface{}) error {
	if mc.conn == nil {
		return fmt.Errorf("not connected")
	}

	return mc.conn.WriteJSON(message)
}

// readMessages reads incoming WebSocket messages
func (mc *MultiplayerClient) readMessages() {
	defer func() {
		mc.connected = false
		if mc.conn != nil {
			mc.conn.Close()
		}
	}()

	for {
		var message map[string]interface{}
		err := mc.conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Multiplayer: Connection error: %v", err)
			break
		}

		// Handle special messages
		if msgType, ok := message["type"].(string); ok {
			switch msgType {
			case "match_found":
				if gameID, ok := message["gameId"].(string); ok {
					mc.gameID = gameID
					log.Printf("Multiplayer: Match found! Game ID: %s", gameID)
				}
			case "rematch_start":
				if gameID, ok := message["gameId"].(string); ok {
					mc.gameID = gameID
					log.Printf("Multiplayer: Rematch started! Game ID: %s", gameID)
				}
			case "game_over":
				mc.gameID = ""
				log.Printf("Multiplayer: Game over")
			}
		}

		// Send to message channel
		select {
		case mc.messages <- message:
		default:
			// Channel full, drop message
			log.Printf("Multiplayer: Message channel full, dropping message")
		}
	}
}
