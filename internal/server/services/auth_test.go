package services

import (
	"testing"

	"github.com/briancain/go-tetris/internal/server/storage/memory"
)

func TestAuthService_Login(t *testing.T) {
	playerStore := memory.NewPlayerStore()
	authService := NewAuthService(playerStore)

	// Test successful login
	player, err := authService.Login("testuser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if player.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", player.Username)
	}

	if player.ID == "" {
		t.Error("Expected player ID to be generated")
	}

	if player.SessionToken == "" {
		t.Error("Expected session token to be generated")
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	playerStore := memory.NewPlayerStore()
	authService := NewAuthService(playerStore)

	// Create a player
	player, err := authService.Login("testuser")
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	// Test valid token
	validatedPlayer, err := authService.ValidateToken(player.SessionToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if validatedPlayer.ID != player.ID {
		t.Errorf("Expected player ID %s, got %s", player.ID, validatedPlayer.ID)
	}

	// Test invalid token
	_, err = authService.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestAuthService_Logout(t *testing.T) {
	playerStore := memory.NewPlayerStore()
	authService := NewAuthService(playerStore)

	// Create a player
	player, err := authService.Login("testuser")
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	// Logout
	err = authService.Logout(player.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify token is invalid after logout
	_, err = authService.ValidateToken(player.SessionToken)
	if err == nil {
		t.Error("Expected error for token after logout")
	}
}
