package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/briancain/go-tetris/internal/server/storage"
	"github.com/briancain/go-tetris/pkg/models"
)

// LeaderboardHandler handles leaderboard-related requests
type LeaderboardHandler struct {
	playerStore storage.PlayerStore
}

// NewLeaderboardHandler creates a new leaderboard handler
func NewLeaderboardHandler(playerStore storage.PlayerStore) *LeaderboardHandler {
	return &LeaderboardHandler{
		playerStore: playerStore,
	}
}

// GetLeaderboard returns the top players by high score
func (h *LeaderboardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	// Get limit from query parameter (default 10)
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get all players (in a real system, we'd have a more efficient query)
	players, err := h.getAllPlayers()
	if err != nil {
		http.Error(w, "Failed to get players", http.StatusInternalServerError)
		return
	}

	// Filter out players who haven't played any games
	var playersWithGames []*models.Player
	for _, player := range players {
		if player.TotalGames > 0 {
			playersWithGames = append(playersWithGames, player)
		}
	}

	// Sort by high score (descending)
	sort.Slice(playersWithGames, func(i, j int) bool {
		return playersWithGames[i].HighScore > playersWithGames[j].HighScore
	})

	// Limit results
	if len(playersWithGames) > limit {
		playersWithGames = playersWithGames[:limit]
	}

	// Create leaderboard response
	leaderboard := make([]LeaderboardEntry, len(playersWithGames))
	for i, player := range playersWithGames {
		leaderboard[i] = LeaderboardEntry{
			Rank:       i + 1,
			Username:   player.Username,
			HighScore:  player.HighScore,
			TotalGames: player.TotalGames,
			Wins:       player.Wins,
			Losses:     player.Losses,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

// LeaderboardEntry represents a leaderboard entry
type LeaderboardEntry struct {
	Rank       int    `json:"rank"`
	Username   string `json:"username"`
	HighScore  int    `json:"highScore"`
	TotalGames int    `json:"totalGames"`
	Wins       int    `json:"wins"`
	Losses     int    `json:"losses"`
}

// getAllPlayers gets all players
func (h *LeaderboardHandler) getAllPlayers() ([]*models.Player, error) {
	return h.playerStore.GetAllPlayers()
}
