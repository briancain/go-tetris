package services

import (
	"math/rand"
	"sync"
	"time"

	"github.com/briancain/go-tetris/internal/server/storage"
	"github.com/briancain/go-tetris/pkg/models"
)

// MatchmakingService handles player matchmaking
type MatchmakingService struct {
	playerStore storage.PlayerStore
	gameStore   storage.GameStore
	queueStore  storage.QueueStore
	gameManager *GameManager
	mu          sync.Mutex // Protects matchmaking operations
}

// NewMatchmakingService creates a new matchmaking service
func NewMatchmakingService(
	playerStore storage.PlayerStore,
	gameStore storage.GameStore,
	queueStore storage.QueueStore,
	gameManager *GameManager,
) *MatchmakingService {
	return &MatchmakingService{
		playerStore: playerStore,
		gameStore:   gameStore,
		queueStore:  queueStore,
		gameManager: gameManager,
	}
}

// JoinQueue adds a player to the matchmaking queue
func (s *MatchmakingService) JoinQueue(playerID string) error {
	// Get player
	player, err := s.playerStore.GetPlayer(playerID)
	if err != nil {
		return err
	}

	// Check if player is already in a game
	if player.GameID != "" {
		return nil // Already in game
	}

	// Add to queue
	err = s.queueStore.AddToQueue(playerID)
	if err != nil {
		return err
	}

	// Update player status
	player.InQueue = true
	err = s.playerStore.UpdatePlayer(player)
	if err != nil {
		return err
	}

	// Try to find a match
	go s.tryMatchmaking()

	return nil
}

// LeaveQueue removes a player from the matchmaking queue
func (s *MatchmakingService) LeaveQueue(playerID string) error {
	// Remove from queue
	err := s.queueStore.RemoveFromQueue(playerID)
	if err != nil {
		return err
	}

	// Update player status
	player, err := s.playerStore.GetPlayer(playerID)
	if err != nil {
		return err
	}

	player.InQueue = false
	return s.playerStore.UpdatePlayer(player)
}

// GetQueueStatus returns the player's position in queue
func (s *MatchmakingService) GetQueueStatus(playerID string) (int, error) {
	return s.queueStore.GetQueuePosition(playerID)
}

// tryMatchmaking attempts to create matches from queued players
func (s *MatchmakingService) tryMatchmaking() {
	s.mu.Lock()
	defer s.mu.Unlock()

	queuedPlayers, err := s.queueStore.GetQueuedPlayers()
	if err != nil || len(queuedPlayers) < 2 {
		return
	}

	// Take first two players
	player1ID := queuedPlayers[0]
	player2ID := queuedPlayers[1]

	// Get player objects
	player1, err := s.playerStore.GetPlayer(player1ID)
	if err != nil {
		return
	}

	player2, err := s.playerStore.GetPlayer(player2ID)
	if err != nil {
		return
	}

	// Remove from queue
	_ = s.queueStore.RemoveFromQueue(player1ID)
	_ = s.queueStore.RemoveFromQueue(player2ID)

	// Create game session
	gameID := generateID()
	seed := generateSeed()

	game := &models.GameSession{
		ID:        gameID,
		Player1:   player1,
		Player2:   player2,
		Seed:      seed,
		Status:    models.GameStatusWaiting,
		CreatedAt: time.Now(),
	}

	// Store game
	err = s.gameStore.CreateGame(game)
	if err != nil {
		// Re-add players to queue on error
		_ = s.queueStore.AddToQueue(player1ID)
		_ = s.queueStore.AddToQueue(player2ID)
		return
	}

	// Update players
	player1.InQueue = false
	player1.GameID = gameID
	player2.InQueue = false
	player2.GameID = gameID

	_ = s.playerStore.UpdatePlayer(player1)
	_ = s.playerStore.UpdatePlayer(player2)

	// Notify game manager
	s.gameManager.StartGame(game)
}

// generateSeed creates a random seed for the game
func generateSeed() int64 {
	return rand.Int63()
}
