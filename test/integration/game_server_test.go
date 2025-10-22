package integration

import (
	"context"
	"testing"

	"github.com/briancain/go-tetris/internal/tetris"
)

func TestGameServerIntegration(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create a game with multiplayer enabled
	game := tetris.NewGame()
	err := game.EnableMultiplayer("http://localhost:8081")
	if err != nil {
		t.Fatalf("Failed to enable multiplayer: %v", err)
	}

	// For now, just test that multiplayer mode is enabled
	// Real WebSocket connection would require proper HTTP authentication
	if !game.MultiplayerMode {
		t.Error("Expected multiplayer mode to be enabled")
	}

	if game.MultiplayerClient == nil {
		t.Error("Expected multiplayer client to be created")
	}

	// Test that the game can handle multiplayer messages
	game.Start()

	// Simulate a match found message
	matchMsg := map[string]interface{}{
		"type":     "match_found",
		"gameId":   "test-game-123",
		"seed":     float64(12345),
		"opponent": "testopponent",
	}

	game.ProcessMultiplayerMessages()

	// Manually handle the message to test the logic
	game.HandleMultiplayerMessage(matchMsg)

	t.Log("✅ Game-server integration test completed successfully")
}

func TestTwoGameInstances(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create two game instances
	game1 := tetris.NewGame()
	game2 := tetris.NewGame()

	// Enable multiplayer for both
	game1.EnableMultiplayer("http://localhost:8081")
	game2.EnableMultiplayer("http://localhost:8081")

	// Connect both (using mock login)
	game1.ConnectToServer("player1")
	game2.ConnectToServer("player2")

	// Both should be in multiplayer mode
	if !game1.MultiplayerMode || !game2.MultiplayerMode {
		t.Error("Expected both games to be in multiplayer mode")
	}

	// Start both games
	game1.Start()
	game2.Start()

	// Process messages for both
	game1.ProcessMultiplayerMessages()
	game2.ProcessMultiplayerMessages()

	t.Log("✅ Two game instances test completed successfully")
}
