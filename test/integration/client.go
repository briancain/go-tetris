package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// TestClient simulates a game client for testing
type TestClient struct {
	PlayerID     string
	Username     string
	SessionToken string
	WSConn       *websocket.Conn
	ServerURL    string
	Messages     []map[string]interface{}
}

// NewTestClient creates a new test client
func NewTestClient(username, serverURL string) *TestClient {
	return &TestClient{
		Username:  username,
		ServerURL: serverURL,
		Messages:  make([]map[string]interface{}, 0),
	}
}

// Login authenticates the client
func (c *TestClient) Login() error {
	loginReq := map[string]string{
		"username": c.Username,
	}

	jsonData, _ := json.Marshal(loginReq)
	resp, err := http.Post(c.ServerURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	var loginResp struct {
		PlayerID     string `json:"playerId"`
		Username     string `json:"username"`
		SessionToken string `json:"sessionToken"`
	}

	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return err
	}

	c.PlayerID = loginResp.PlayerID
	c.SessionToken = loginResp.SessionToken

	return nil
}

// ConnectWebSocket establishes WebSocket connection
func (c *TestClient) ConnectWebSocket() error {
	u, _ := url.Parse(c.ServerURL)
	u.Scheme = "ws"
	u.Path = "/ws"
	u.RawQuery = "token=" + c.SessionToken

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.WSConn = conn

	// Start message reader
	go c.readMessages()

	return nil
}

// JoinQueue joins the matchmaking queue
func (c *TestClient) JoinQueue() error {
	req, _ := http.NewRequest("POST", c.ServerURL+"/api/matchmaking/queue", nil)
	req.Header.Set("Authorization", "Bearer "+c.SessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("join queue failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetQueueStatus checks queue position
func (c *TestClient) GetQueueStatus() (int, error) {
	req, _ := http.NewRequest("GET", c.ServerURL+"/api/matchmaking/status", nil)
	req.Header.Set("Authorization", "Bearer "+c.SessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var statusResp struct {
		Position int  `json:"position"`
		InQueue  bool `json:"inQueue"`
	}

	err = json.NewDecoder(resp.Body).Decode(&statusResp)
	if err != nil {
		return -1, err
	}

	return statusResp.Position, nil
}

// SendGameMove sends a game move via WebSocket
func (c *TestClient) SendGameMove(moveType string) error {
	if c.WSConn == nil {
		return fmt.Errorf("WebSocket not connected")
	}

	message := map[string]interface{}{
		"type":     "game_move",
		"moveType": moveType,
	}

	return c.WSConn.WriteJSON(message)
}

// SendGameState sends game state via WebSocket
func (c *TestClient) SendGameState(board [][]int, score, level, lines int) error {
	if c.WSConn == nil {
		return fmt.Errorf("WebSocket not connected")
	}

	message := map[string]interface{}{
		"type":  "game_state",
		"board": board,
		"score": score,
		"level": level,
		"lines": lines,
	}

	return c.WSConn.WriteJSON(message)
}

// SendPing sends a ping message
func (c *TestClient) SendPing() error {
	if c.WSConn == nil {
		return fmt.Errorf("WebSocket not connected")
	}

	message := map[string]interface{}{
		"type": "ping",
	}

	return c.WSConn.WriteJSON(message)
}

// WaitForMessage waits for a specific message type
func (c *TestClient) WaitForMessage(messageType string, timeout time.Duration) (map[string]interface{}, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		for _, msg := range c.Messages {
			if msgType, ok := msg["type"].(string); ok && msgType == messageType {
				return msg, nil
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	return nil, fmt.Errorf("timeout waiting for message type: %s", messageType)
}

// GetMessages returns all received messages
func (c *TestClient) GetMessages() []map[string]interface{} {
	return c.Messages
}

// Close closes the WebSocket connection
func (c *TestClient) Close() {
	if c.WSConn != nil {
		c.WSConn.Close()
	}
}

// readMessages reads incoming WebSocket messages
func (c *TestClient) readMessages() {
	for {
		var message map[string]interface{}
		err := c.WSConn.ReadJSON(&message)
		if err != nil {
			break
		}

		c.Messages = append(c.Messages, message)
	}
}
