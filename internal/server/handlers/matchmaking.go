package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/briancain/go-tetris/internal/server/logger"
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
	requestID := r.Context().Value("requestID").(string)

	if r.Method != http.MethodPost {
		logger.Logger.Warn("Invalid method for join queue",
			"requestID", requestID,
			"method", r.Method,
		)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get player from context (set by auth middleware)
	playerID := r.Context().Value("playerID").(string)

	err := h.matchmakingService.JoinQueue(playerID)
	if err != nil {
		logger.Logger.Error("Failed to join matchmaking queue",
			"requestID", requestID,
			"playerID", playerID,
			"error", err,
		)
		http.Error(w, "Failed to join queue", http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Player joined matchmaking queue",
		"requestID", requestID,
		"playerID", playerID,
	)

	w.WriteHeader(http.StatusOK)
}

// LeaveQueue handles leaving the matchmaking queue
func (h *MatchmakingHandler) LeaveQueue(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID").(string)

	if r.Method != http.MethodDelete {
		logger.Logger.Warn("Invalid method for leave queue",
			"requestID", requestID,
			"method", r.Method,
		)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get player from context (set by auth middleware)
	playerID := r.Context().Value("playerID").(string)

	err := h.matchmakingService.LeaveQueue(playerID)
	if err != nil {
		logger.Logger.Error("Failed to leave matchmaking queue",
			"requestID", requestID,
			"playerID", playerID,
			"error", err,
		)
		http.Error(w, "Failed to leave queue", http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Player left matchmaking queue",
		"requestID", requestID,
		"playerID", playerID,
	)

	w.WriteHeader(http.StatusOK)
}

// GetQueueStatus handles getting queue status
func (h *MatchmakingHandler) GetQueueStatus(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID").(string)

	if r.Method != http.MethodGet {
		logger.Logger.Warn("Invalid method for queue status",
			"requestID", requestID,
			"method", r.Method,
		)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get player from context (set by auth middleware)
	playerID := r.Context().Value("playerID").(string)

	position, err := h.matchmakingService.GetQueueStatus(playerID)
	if err != nil {
		logger.Logger.Error("Failed to get queue status",
			"requestID", requestID,
			"playerID", playerID,
			"error", err,
		)
		http.Error(w, "Failed to get queue status", http.StatusInternalServerError)
		return
	}

	response := QueueStatusResponse{
		Position: position,
		InQueue:  position >= 0,
	}

	logger.Logger.Debug("Queue status retrieved",
		"requestID", requestID,
		"playerID", playerID,
		"position", position,
		"inQueue", response.InQueue,
	)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger.Error("Failed to encode queue status response",
			"requestID", requestID,
			"playerID", playerID,
			"error", err,
		)
	}
}
