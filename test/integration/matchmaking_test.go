package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/briancain/go-tetris/internal/server/handlers"
	"github.com/briancain/go-tetris/internal/server/middleware"
	"github.com/briancain/go-tetris/internal/server/services"
	"github.com/briancain/go-tetris/internal/server/storage/memory"
)

const testServerURL = "http://localhost:8081"

// startTestServer starts a test server on port 8081
func startTestServer() *http.Server {
	// Initialize storage
	playerStore := memory.NewPlayerStore()
	gameStore := memory.NewGameStore()
	queueStore := memory.NewQueueStore()

	// Initialize services
	authService := services.NewAuthService(playerStore)
	wsManager := services.NewWebSocketManager()
	gameManager := services.NewGameManager(gameStore, playerStore, wsManager)
	matchmakingService := services.NewMatchmakingService(playerStore, gameStore, queueStore, gameManager)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	matchmakingHandler := handlers.NewMatchmakingHandler(matchmakingService)
	wsHandler := handlers.NewWebSocketHandler(wsManager, authService, gameManager)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/logout", authMiddleware.RequireAuth(authHandler.Logout))
	mux.HandleFunc("/api/matchmaking/queue", authMiddleware.RequireAuth(matchmakingHandler.JoinQueue))
	mux.HandleFunc("/api/matchmaking/queue/leave", authMiddleware.RequireAuth(matchmakingHandler.LeaveQueue))
	mux.HandleFunc("/api/matchmaking/status", authMiddleware.RequireAuth(matchmakingHandler.GetQueueStatus))
	mux.HandleFunc("/ws", wsHandler.HandleWebSocket)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go server.ListenAndServe()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	return server
}

func TestFullMatchmakingFlow(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create two test clients
	client1 := NewTestClient("player1", testServerURL)
	client2 := NewTestClient("player2", testServerURL)

	// Step 1: Both clients login
	err := client1.Login()
	if err != nil {
		t.Fatalf("Client1 login failed: %v", err)
	}

	err = client2.Login()
	if err != nil {
		t.Fatalf("Client2 login failed: %v", err)
	}

	// Step 2: Connect WebSockets
	err = client1.ConnectWebSocket()
	if err != nil {
		t.Fatalf("Client1 WebSocket connection failed: %v", err)
	}
	defer client1.Close()

	err = client2.ConnectWebSocket()
	if err != nil {
		t.Fatalf("Client2 WebSocket connection failed: %v", err)
	}
	defer client2.Close()

	// Step 3: Both join queue
	err = client1.JoinQueue()
	if err != nil {
		t.Fatalf("Client1 join queue failed: %v", err)
	}

	err = client2.JoinQueue()
	if err != nil {
		t.Fatalf("Client2 join queue failed: %v", err)
	}

	// Step 4: Wait for match found messages
	matchMsg1, err := client1.WaitForMessage("match_found", 2*time.Second)
	if err != nil {
		t.Fatalf("Client1 didn't receive match_found: %v", err)
	}

	matchMsg2, err := client2.WaitForMessage("match_found", 2*time.Second)
	if err != nil {
		t.Fatalf("Client2 didn't receive match_found: %v", err)
	}

	// Step 5: Verify match details
	gameID1, ok := matchMsg1["gameId"].(string)
	if !ok || gameID1 == "" {
		t.Error("Client1 match message missing gameId")
	}

	gameID2, ok := matchMsg2["gameId"].(string)
	if !ok || gameID2 == "" {
		t.Error("Client2 match message missing gameId")
	}

	if gameID1 != gameID2 {
		t.Errorf("Clients got different game IDs: %s vs %s", gameID1, gameID2)
	}

	// Verify seeds are the same
	seed1, ok := matchMsg1["seed"].(float64)
	if !ok {
		t.Error("Client1 match message missing seed")
	}

	seed2, ok := matchMsg2["seed"].(float64)
	if !ok {
		t.Error("Client2 match message missing seed")
	}

	if seed1 != seed2 {
		t.Errorf("Clients got different seeds: %f vs %f", seed1, seed2)
	}

	// Verify opponent info
	opponent1, ok := matchMsg1["opponent"].(string)
	if !ok || opponent1 != "player2" {
		t.Errorf("Client1 expected opponent 'player2', got '%s'", opponent1)
	}

	opponent2, ok := matchMsg2["opponent"].(string)
	if !ok || opponent2 != "player1" {
		t.Errorf("Client2 expected opponent 'player1', got '%s'", opponent2)
	}

	t.Logf("✅ Match created successfully: GameID=%s, Seed=%.0f", gameID1, seed1)
}

func TestGameMoveExchange(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create and setup two clients
	client1 := NewTestClient("player1", testServerURL)
	client2 := NewTestClient("player2", testServerURL)

	// Login and connect
	client1.Login()
	client2.Login()
	client1.ConnectWebSocket()
	client2.ConnectWebSocket()
	defer client1.Close()
	defer client2.Close()

	// Join queue and wait for match
	client1.JoinQueue()
	client2.JoinQueue()

	client1.WaitForMessage("match_found", 2*time.Second)
	client2.WaitForMessage("match_found", 2*time.Second)

	// Clear previous messages
	client1.Messages = nil
	client2.Messages = nil

	// Step 1: Client1 sends a move
	err := client1.SendGameMove("left")
	if err != nil {
		t.Fatalf("Failed to send move: %v", err)
	}

	// Step 2: Client2 should receive the move
	moveMsg, err := client2.WaitForMessage("game_move", 1*time.Second)
	if err != nil {
		t.Fatalf("Client2 didn't receive move: %v", err)
	}

	// Verify move details
	moveType, ok := moveMsg["moveType"].(string)
	if !ok || moveType != "left" {
		t.Errorf("Expected moveType 'left', got '%s'", moveType)
	}

	playerID, ok := moveMsg["playerId"].(string)
	if !ok || playerID != client1.PlayerID {
		t.Errorf("Expected playerId '%s', got '%s'", client1.PlayerID, playerID)
	}

	t.Logf("✅ Move exchange successful: %s sent 'left' move to %s", client1.Username, client2.Username)
}

func TestGameStateSync(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create and setup two clients
	client1 := NewTestClient("player1", testServerURL)
	client2 := NewTestClient("player2", testServerURL)

	// Login and connect
	client1.Login()
	client2.Login()
	client1.ConnectWebSocket()
	client2.ConnectWebSocket()
	defer client1.Close()
	defer client2.Close()

	// Join queue and wait for match
	client1.JoinQueue()
	client2.JoinQueue()

	client1.WaitForMessage("match_found", 2*time.Second)
	client2.WaitForMessage("match_found", 2*time.Second)

	// Clear previous messages
	client1.Messages = nil
	client2.Messages = nil

	// Step 1: Client1 sends game state
	testBoard := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0, 0, 1, 1},
	}

	err := client1.SendGameState(testBoard, 1200, 3, 5)
	if err != nil {
		t.Fatalf("Failed to send game state: %v", err)
	}

	// Step 2: Client2 should receive the state
	stateMsg, err := client2.WaitForMessage("game_state", 1*time.Second)
	if err != nil {
		t.Fatalf("Client2 didn't receive game state: %v", err)
	}

	// Verify state details
	score, ok := stateMsg["score"].(float64)
	if !ok || int(score) != 1200 {
		t.Errorf("Expected score 1200, got %v", score)
	}

	level, ok := stateMsg["level"].(float64)
	if !ok || int(level) != 3 {
		t.Errorf("Expected level 3, got %v", level)
	}

	lines, ok := stateMsg["lines"].(float64)
	if !ok || int(lines) != 5 {
		t.Errorf("Expected lines 5, got %v", lines)
	}

	t.Logf("✅ Game state sync successful: Score=%d, Level=%d, Lines=%d",
		int(score), int(level), int(lines))
}

func TestPingPong(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create and setup client
	client := NewTestClient("player1", testServerURL)
	client.Login()
	client.ConnectWebSocket()
	defer client.Close()

	// Send ping
	err := client.SendPing()
	if err != nil {
		t.Fatalf("Failed to send ping: %v", err)
	}

	// Wait for pong
	pongMsg, err := client.WaitForMessage("pong", 1*time.Second)
	if err != nil {
		t.Fatalf("Didn't receive pong: %v", err)
	}

	msgType, ok := pongMsg["type"].(string)
	if !ok || msgType != "pong" {
		t.Errorf("Expected pong message, got %v", pongMsg)
	}

	t.Logf("✅ Ping-pong successful")
}
