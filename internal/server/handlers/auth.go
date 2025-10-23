package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/briancain/go-tetris/internal/server/logger"
	"github.com/briancain/go-tetris/internal/server/services"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	PlayerID     string `json:"playerId"`
	Username     string `json:"username"`
	SessionToken string `json:"sessionToken"`
}

// Login handles player login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID").(string)

	if r.Method != http.MethodPost {
		logger.Logger.Warn("Invalid method for login",
			"requestID", requestID,
			"method", r.Method,
		)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Logger.Error("Failed to decode login request",
			"requestID", requestID,
			"error", err,
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		logger.Logger.Warn("Login attempt with empty username",
			"requestID", requestID,
		)
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Create player session
	player, err := h.authService.Login(req.Username)
	if err != nil {
		// Check if it's a username conflict error
		if err.Error() == fmt.Sprintf("username '%s' is already in use", req.Username) {
			logger.Logger.Warn("Login attempt with username already in use",
				"requestID", requestID,
				"username", req.Username,
			)
			http.Error(w, "Username is already in use. Please choose a different username.", http.StatusConflict)
			return
		}

		logger.Logger.Error("Failed to create player session",
			"requestID", requestID,
			"username", req.Username,
			"error", err,
		)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Player logged in successfully",
		"requestID", requestID,
		"playerID", player.ID,
		"username", player.Username,
	)

	// Return response
	response := LoginResponse{
		PlayerID:     player.ID,
		Username:     player.Username,
		SessionToken: player.SessionToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger.Error("Failed to encode login response",
			"requestID", requestID,
			"playerID", player.ID,
			"error", err,
		)
	}
}

// Logout handles player logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID").(string)

	if r.Method != http.MethodPost {
		logger.Logger.Warn("Invalid method for logout",
			"requestID", requestID,
			"method", r.Method,
		)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get player from context (set by auth middleware)
	playerID := r.Context().Value("playerID").(string)

	err := h.authService.Logout(playerID)
	if err != nil {
		logger.Logger.Error("Failed to logout player",
			"requestID", requestID,
			"playerID", playerID,
			"error", err,
		)
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Player logged out successfully",
		"requestID", requestID,
		"playerID", playerID,
	)

	w.WriteHeader(http.StatusOK)
}
