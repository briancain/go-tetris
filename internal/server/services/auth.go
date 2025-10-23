package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/briancain/go-tetris/internal/server/storage"
	"github.com/briancain/go-tetris/pkg/models"
)

var ErrUsernameInUse = errors.New("username is already in use")

// AuthService handles player authentication
type AuthService struct {
	playerStore storage.PlayerStore
}

// NewAuthService creates a new authentication service
func NewAuthService(playerStore storage.PlayerStore) *AuthService {
	return &AuthService{
		playerStore: playerStore,
	}
}

// Login creates a new player session
func (s *AuthService) Login(username string) (*models.Player, error) {
	// Check if username is already taken by an active player
	_, err := s.playerStore.GetPlayerByUsername(username)
	if err == nil {
		// Username exists, check if player is still active
		// For now, we'll consider any existing player as active
		// In a production system, you might want to check last activity time
		return nil, ErrUsernameInUse
	}

	// Generate unique player ID and session token
	playerID := generateID()
	sessionToken := generateSessionToken()

	player := &models.Player{
		ID:           playerID,
		Username:     username,
		SessionToken: sessionToken,
		ConnectedAt:  time.Now(),
		LastActivity: time.Now(),
		InQueue:      false,
	}

	err = s.playerStore.CreatePlayer(player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

// ValidateToken checks if a session token is valid and returns the player
func (s *AuthService) ValidateToken(token string) (*models.Player, error) {
	return s.playerStore.GetPlayerByToken(token)
}

// GetPlayerByID retrieves a player by their ID
func (s *AuthService) GetPlayerByID(playerID string) (*models.Player, error) {
	return s.playerStore.GetPlayer(playerID)
}

// DeletePlayer removes a player session
func (s *AuthService) DeletePlayer(playerID string) error {
	return s.playerStore.DeletePlayer(playerID)
}

// Logout removes a player session
func (s *AuthService) Logout(playerID string) error {
	return s.playerStore.DeletePlayer(playerID)
}

// UpdatePlayerActivity updates the last activity time for a player
func (s *AuthService) UpdatePlayerActivity(playerID string) error {
	player, err := s.playerStore.GetPlayer(playerID)
	if err != nil {
		return err
	}

	player.LastActivity = time.Now()
	return s.playerStore.UpdatePlayer(player)
}

// CleanupInactivePlayers removes players who have been inactive for too long
func (s *AuthService) CleanupInactivePlayers(inactiveThreshold time.Duration) error {
	players, err := s.playerStore.GetAllPlayers()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, player := range players {
		if now.Sub(player.LastActivity) > inactiveThreshold {
			// Remove inactive player to free up username
			_ = s.playerStore.DeletePlayer(player.ID)
		}
	}

	return nil
}

// generateID creates a unique identifier
func generateID() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateSessionToken creates a session token
func generateSessionToken() string {
	bytes := make([]byte, 32)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
