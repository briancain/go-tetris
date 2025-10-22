package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/briancain/go-tetris/internal/server/storage/memory"
	"github.com/briancain/go-tetris/pkg/models"
)

func TestGetLeaderboard_FiltersPlayersWithNoGames(t *testing.T) {
	// Setup
	playerStore := memory.NewPlayerStore()
	handler := NewLeaderboardHandler(playerStore)

	// Create players - some with games, some without
	playerWithGames := &models.Player{
		ID:         "player1",
		Username:   "HasPlayed",
		TotalGames: 5,
		Wins:       3,
		Losses:     2,
		HighScore:  1500,
	}

	playerWithoutGames := &models.Player{
		ID:         "player2", 
		Username:   "JustLoggedIn",
		TotalGames: 0, // No games played
		Wins:       0,
		Losses:     0,
		HighScore:  0,
	}

	anotherPlayerWithGames := &models.Player{
		ID:         "player3",
		Username:   "AlsoPlayed", 
		TotalGames: 2,
		Wins:       1,
		Losses:     1,
		HighScore:  800,
	}

	playerStore.CreatePlayer(playerWithGames)
	playerStore.CreatePlayer(playerWithoutGames)
	playerStore.CreatePlayer(anotherPlayerWithGames)

	// Make request
	req := httptest.NewRequest("GET", "/api/leaderboard", nil)
	w := httptest.NewRecorder()

	handler.GetLeaderboard(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var leaderboard []LeaderboardEntry
	err := json.NewDecoder(w.Body).Decode(&leaderboard)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should only have 2 players (those who played games)
	if len(leaderboard) != 2 {
		t.Errorf("Expected 2 players in leaderboard, got %d", len(leaderboard))
	}

	// Should be sorted by high score (descending)
	if leaderboard[0].Username != "HasPlayed" || leaderboard[0].HighScore != 1500 {
		t.Errorf("Expected first player to be HasPlayed with score 1500, got %s with %d", 
			leaderboard[0].Username, leaderboard[0].HighScore)
	}

	if leaderboard[1].Username != "AlsoPlayed" || leaderboard[1].HighScore != 800 {
		t.Errorf("Expected second player to be AlsoPlayed with score 800, got %s with %d",
			leaderboard[1].Username, leaderboard[1].HighScore)
	}

	// Verify player without games is not included
	for _, entry := range leaderboard {
		if entry.Username == "JustLoggedIn" {
			t.Error("Player with 0 games should not appear in leaderboard")
		}
	}
}

func TestGetLeaderboard_EmptyWhenNoPlayersHaveGames(t *testing.T) {
	// Setup
	playerStore := memory.NewPlayerStore()
	handler := NewLeaderboardHandler(playerStore)

	// Create only players without games
	player1 := &models.Player{
		ID:         "player1",
		Username:   "NoGames1",
		TotalGames: 0,
		HighScore:  0,
	}

	player2 := &models.Player{
		ID:         "player2",
		Username:   "NoGames2", 
		TotalGames: 0,
		HighScore:  0,
	}

	playerStore.CreatePlayer(player1)
	playerStore.CreatePlayer(player2)

	// Make request
	req := httptest.NewRequest("GET", "/api/leaderboard", nil)
	w := httptest.NewRecorder()

	handler.GetLeaderboard(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var leaderboard []LeaderboardEntry
	err := json.NewDecoder(w.Body).Decode(&leaderboard)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should be empty
	if len(leaderboard) != 0 {
		t.Errorf("Expected empty leaderboard, got %d entries", len(leaderboard))
	}
}
