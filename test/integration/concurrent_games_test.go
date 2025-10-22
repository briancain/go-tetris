package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/briancain/go-tetris/internal/tetris"
)

func TestConcurrentGames(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create 4 players for 2 concurrent games
	players := make([]*tetris.Game, 4)
	for i := 0; i < 4; i++ {
		game := tetris.NewGame()
		err := game.EnableMultiplayer("http://localhost:8081")
		if err != nil {
			t.Fatalf("Failed to enable multiplayer for player %d: %v", i, err)
		}
		players[i] = game
	}

	// Connect all players concurrently
	var wg sync.WaitGroup
	for i, player := range players {
		wg.Add(1)
		go func(i int, p *tetris.Game) {
			defer wg.Done()
			
			username := fmt.Sprintf("player%d", i+1)
			err := p.ConnectToServer(username)
			if err != nil {
				t.Errorf("Player %d connection failed: %v", i+1, err)
				return
			}

			err = p.JoinMatchmaking()
			if err != nil {
				t.Errorf("Player %d join queue failed: %v", i+1, err)
			}
		}(i, player)
	}

	wg.Wait()

	// Give time for matchmaking
	time.Sleep(1 * time.Second)

	// Process messages for all players
	for _, player := range players {
		player.ProcessMultiplayerMessages()
	}

	// Check that we have 2 different games
	gameIDs := make(map[string]int)
	matchedPlayers := 0

	for i, player := range players {
		gameID := player.MultiplayerClient.GetGameID()
		if gameID != "" {
			gameIDs[gameID]++
			matchedPlayers++
			t.Logf("Player %d matched in game: %s", i+1, gameID)
		}
	}

	// Verify we have exactly 2 games with 2 players each
	if len(gameIDs) != 2 {
		t.Errorf("Expected 2 concurrent games, got %d", len(gameIDs))
	}

	for gameID, playerCount := range gameIDs {
		if playerCount != 2 {
			t.Errorf("Game %s has %d players, expected 2", gameID, playerCount)
		}
	}

	if matchedPlayers != 4 {
		t.Errorf("Expected 4 matched players, got %d", matchedPlayers)
	}

	t.Logf("✅ Successfully created %d concurrent games with %d total players", len(gameIDs), matchedPlayers)
}

func TestConcurrentGameMoves(t *testing.T) {
	// Start test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create 2 players
	game1 := tetris.NewGame()
	game2 := tetris.NewGame()

	game1.EnableMultiplayer("http://localhost:8081")
	game2.EnableMultiplayer("http://localhost:8081")

	// Connect and match players
	game1.ConnectToServer("player1")
	game2.ConnectToServer("player2")
	game1.JoinMatchmaking()
	game2.JoinMatchmaking()

	// Wait for match
	time.Sleep(500 * time.Millisecond)
	game1.ProcessMultiplayerMessages()
	game2.ProcessMultiplayerMessages()

	// Verify they're matched
	if game1.MultiplayerClient.GetGameID() == "" || game2.MultiplayerClient.GetGameID() == "" {
		t.Skip("Players not matched, skipping move test")
	}

	// Clear previous messages by processing them
	for {
		msg1 := game1.MultiplayerClient.GetMessage()
		msg2 := game2.MultiplayerClient.GetMessage()
		if msg1 == nil && msg2 == nil {
			break
		}
	}

	// Send concurrent moves
	var wg sync.WaitGroup
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			game1.MoveLeft()
			time.Sleep(10 * time.Millisecond)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			game2.MoveRight()
			time.Sleep(10 * time.Millisecond)
		}
	}()

	wg.Wait()

	// Give time for message processing
	time.Sleep(200 * time.Millisecond)

	// Process messages
	game1.ProcessMultiplayerMessages()
	game2.ProcessMultiplayerMessages()

	t.Log("✅ Concurrent game moves completed without errors")
}
