package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/briancain/go-tetris/internal/server/storage/memory"
)

func TestAuthService_UsernameConflicts(t *testing.T) {
	playerStore := memory.NewPlayerStore()
	authService := NewAuthService(playerStore)

	t.Run("should allow unique usernames", func(t *testing.T) {
		// First player should succeed
		player1, err := authService.Login("player1")
		if err != nil {
			t.Fatalf("Expected first login to succeed, got error: %v", err)
		}
		if player1.Username != "player1" {
			t.Errorf("Expected username 'player1', got '%s'", player1.Username)
		}

		// Second player with different username should succeed
		player2, err := authService.Login("player2")
		if err != nil {
			t.Fatalf("Expected second login with different username to succeed, got error: %v", err)
		}
		if player2.Username != "player2" {
			t.Errorf("Expected username 'player2', got '%s'", player2.Username)
		}
	})

	t.Run("should reject duplicate usernames", func(t *testing.T) {
		// First player logs in
		_, err := authService.Login("duplicate_user")
		if err != nil {
			t.Fatalf("Expected first login to succeed, got error: %v", err)
		}

		// Second player tries same username
		_, err = authService.Login("duplicate_user")
		if err == nil {
			t.Fatal("Expected second login with same username to fail")
		}

		if !errors.Is(err, ErrUsernameInUse) {
			t.Errorf("Expected ErrUsernameInUse, got '%v'", err)
		}
	})

	t.Run("should allow username reuse after logout", func(t *testing.T) {
		// First player logs in
		player1, err := authService.Login("reusable_user")
		if err != nil {
			t.Fatalf("Expected first login to succeed, got error: %v", err)
		}

		// First player logs out
		err = authService.Logout(player1.ID)
		if err != nil {
			t.Fatalf("Expected logout to succeed, got error: %v", err)
		}

		// Second player should now be able to use the same username
		player2, err := authService.Login("reusable_user")
		if err != nil {
			t.Fatalf("Expected login after logout to succeed, got error: %v", err)
		}
		if player2.Username != "reusable_user" {
			t.Errorf("Expected username 'reusable_user', got '%s'", player2.Username)
		}
		if player1.ID == player2.ID {
			t.Error("Expected different player IDs for different sessions")
		}
	})

	t.Run("should handle case sensitivity", func(t *testing.T) {
		// First player logs in with lowercase
		_, err := authService.Login("testuser")
		if err != nil {
			t.Fatalf("Expected first login to succeed, got error: %v", err)
		}

		// Second player tries with different case - should be treated as different username
		_, err = authService.Login("TestUser")
		if err != nil {
			t.Fatalf("Expected login with different case to succeed, got error: %v", err)
		}
	})
}

func TestAuthService_CleanupInactivePlayers(t *testing.T) {
	playerStore := memory.NewPlayerStore()
	authService := NewAuthService(playerStore)

	t.Run("should cleanup inactive players", func(t *testing.T) {
		// Create a player
		player, err := authService.Login("inactive_user")
		if err != nil {
			t.Fatalf("Expected login to succeed, got error: %v", err)
		}

		// Manually set last activity time to past (simulate inactive player)
		player.LastActivity = time.Now().Add(-2 * time.Hour)
		err = playerStore.UpdatePlayer(player)
		if err != nil {
			t.Fatalf("Failed to update player: %v", err)
		}

		// Run cleanup with 1 hour threshold
		err = authService.CleanupInactivePlayers(1 * time.Hour)
		if err != nil {
			t.Fatalf("Expected cleanup to succeed, got error: %v", err)
		}

		// Player should be removed, so username should be available again
		_, err = authService.Login("inactive_user")
		if err != nil {
			t.Fatalf("Expected login after cleanup to succeed, got error: %v", err)
		}
	})

	t.Run("should not cleanup active players", func(t *testing.T) {
		// Create a player
		_, err := authService.Login("active_user")
		if err != nil {
			t.Fatalf("Expected login to succeed, got error: %v", err)
		}

		// Run cleanup with 1 hour threshold (player is recent)
		err = authService.CleanupInactivePlayers(1 * time.Hour)
		if err != nil {
			t.Fatalf("Expected cleanup to succeed, got error: %v", err)
		}

		// Player should still exist, so username should be taken
		_, err = authService.Login("active_user")
		if err == nil {
			t.Fatal("Expected login with existing username to fail")
		}
		if !strings.Contains(err.Error(), "already in use") {
			t.Errorf("Expected 'already in use' error, got: %v", err)
		}
	})

	t.Run("should allow username reuse after disconnect cleanup", func(t *testing.T) {
		// Player logs in
		player, err := authService.Login("disconnect_user")
		if err != nil {
			t.Fatalf("Expected login to succeed, got error: %v", err)
		}

		// Simulate disconnect by deleting player
		err = authService.DeletePlayer(player.ID)
		if err != nil {
			t.Fatalf("Expected delete to succeed, got error: %v", err)
		}

		// Same username should now be available
		player2, err := authService.Login("disconnect_user")
		if err != nil {
			t.Fatalf("Expected login with same username after disconnect to succeed, got error: %v", err)
		}

		if player2.Username != "disconnect_user" {
			t.Errorf("Expected username 'disconnect_user', got '%s'", player2.Username)
		}

		if player2.ID == player.ID {
			t.Error("Expected new player to have different ID")
		}
	})
}
