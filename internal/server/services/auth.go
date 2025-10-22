package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/briancain/go-tetris/internal/server/storage"
	"github.com/briancain/go-tetris/pkg/models"
)

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
	// Generate unique player ID and session token
	playerID := generateID()
	sessionToken := generateSessionToken()

	player := &models.Player{
		ID:           playerID,
		Username:     username,
		SessionToken: sessionToken,
		ConnectedAt:  time.Now(),
		InQueue:      false,
	}

	err := s.playerStore.CreatePlayer(player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

// ValidateToken checks if a session token is valid and returns the player
func (s *AuthService) ValidateToken(token string) (*models.Player, error) {
	return s.playerStore.GetPlayerByToken(token)
}

// Logout removes a player session
func (s *AuthService) Logout(playerID string) error {
	return s.playerStore.DeletePlayer(playerID)
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
