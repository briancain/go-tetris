package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/briancain/go-tetris/internal/server/services"
)

// MatchmakingHandler handles matchmaking endpoints
type MatchmakingHandler struct {
	matchmakingService *services.MatchmakingService
}

// NewMatchmakingHandler creates a new matchmaking handler
func NewMatchmakingHandler(matchmakingService *services.MatchmakingService) *MatchmakingHandler {
	return &MatchmakingHandler{
		matchmakingService: matchmakingService,
	}
}

// QueueStatusResponse represents queue status response
type QueueStatusResponse struct {
	Position int  `json:"position"`
	InQueue  bool `json:"inQueue"`
}

// JoinQueue handles joining the matchmaking queue
func (h *MatchmakingHandler) JoinQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get player from context (set by auth middleware)
	playerID := r.Context().Value("playerID").(string)

	err := h.matchmakingService.JoinQueue(playerID)
	if err != nil {
		http.Error(w, "Failed to join queue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// LeaveQueue handles leaving the matchmaking queue
func (h *MatchmakingHandler) LeaveQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get player from context (set by auth middleware)
	playerID := r.Context().Value("playerID").(string)

	err := h.matchmakingService.LeaveQueue(playerID)
	if err != nil {
		http.Error(w, "Failed to leave queue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetQueueStatus handles getting queue status
func (h *MatchmakingHandler) GetQueueStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get player from context (set by auth middleware)
	playerID := r.Context().Value("playerID").(string)

	position, err := h.matchmakingService.GetQueueStatus(playerID)
	if err != nil {
		http.Error(w, "Failed to get queue status", http.StatusInternalServerError)
		return
	}

	response := QueueStatusResponse{
		Position: position,
		InQueue:  position >= 0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
