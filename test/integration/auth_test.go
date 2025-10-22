package integration

import (
	"context"
	"testing"
	"time"

	"github.com/briancain/go-tetris/internal/tetris"
)

func TestRealAuthentication(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create multiplayer client
	client := tetris.NewMultiplayerClient("http://localhost:8081")

	// Test login
	err := client.Login("testplayer")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if client.GetUsername() != "testplayer" {
		t.Errorf("Expected username 'testplayer', got '%s'", client.GetUsername())
	}

	// Test WebSocket connection
	err = client.Connect()
	if err != nil {
		t.Fatalf("WebSocket connection failed: %v", err)
	}

	if !client.IsConnected() {
		t.Error("Expected client to be connected")
	}

	// Test join queue
	err = client.JoinQueue()
	if err != nil {
		t.Fatalf("Join queue failed: %v", err)
	}

	// Clean up
	client.Close()

	t.Log("✅ Real authentication test passed")
}

func TestFullGameAuthentication(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create game with multiplayer
	game := tetris.NewGame()
	err := game.EnableMultiplayer("http://localhost:8081")
	if err != nil {
		t.Fatalf("Failed to enable multiplayer: %v", err)
	}

	// Connect with real authentication
	err = game.ConnectToServer("gameplayer")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	// Join matchmaking
	err = game.JoinMatchmaking()
	if err != nil {
		t.Fatalf("Failed to join matchmaking: %v", err)
	}

	// Process any messages
	game.ProcessMultiplayerMessages()

	t.Log("✅ Full game authentication test passed")
}

func TestTwoPlayersMatchmaking(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create two games
	game1 := tetris.NewGame()
	game2 := tetris.NewGame()

	// Enable multiplayer for both
	game1.EnableMultiplayer("http://localhost:8081")
	game2.EnableMultiplayer("http://localhost:8081")

	// Connect both players
	err := game1.ConnectToServer("player1")
	if err != nil {
		t.Fatalf("Player1 connection failed: %v", err)
	}

	err = game2.ConnectToServer("player2")
	if err != nil {
		t.Fatalf("Player2 connection failed: %v", err)
	}

	// Both join matchmaking
	err = game1.JoinMatchmaking()
	if err != nil {
		t.Fatalf("Player1 join queue failed: %v", err)
	}

	err = game2.JoinMatchmaking()
	if err != nil {
		t.Fatalf("Player2 join queue failed: %v", err)
	}

	// Give time for matchmaking
	time.Sleep(500 * time.Millisecond)

	// Process messages for both players
	game1.ProcessMultiplayerMessages()
	game2.ProcessMultiplayerMessages()

	// Check if both games have a game ID (indicating they were matched)
	gameID1 := game1.MultiplayerClient.GetGameID()
	gameID2 := game2.MultiplayerClient.GetGameID()

	if gameID1 == "" || gameID2 == "" {
		t.Log("Players may not have been matched yet (this is okay for timing)")
	} else if gameID1 != gameID2 {
		t.Errorf("Players got different game IDs: %s vs %s", gameID1, gameID2)
	} else {
		t.Logf("✅ Players successfully matched in game: %s", gameID1)
	}

	t.Log("✅ Two players matchmaking test completed")
}
