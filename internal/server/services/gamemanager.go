package services

import (
	"encoding/json"
	"log"

	"github.com/briancain/go-tetris/internal/server/storage"
	"github.com/briancain/go-tetris/pkg/models"
)

// GameManager handles active game sessions
type GameManager struct {
	gameStore   storage.GameStore
	playerStore storage.PlayerStore
	wsManager   *WebSocketManager
}

// NewGameManager creates a new game manager
func NewGameManager(
	gameStore storage.GameStore,
	playerStore storage.PlayerStore,
	wsManager *WebSocketManager,
) *GameManager {
	return &GameManager{
		gameStore:   gameStore,
		playerStore: playerStore,
		wsManager:   wsManager,
	}
}

// StartGame initializes a new game session
func (gm *GameManager) StartGame(game *models.GameSession) {
	// Update game status
	game.Status = models.GameStatusActive
	err := gm.gameStore.UpdateGame(game)
	if err != nil {
		log.Printf("Failed to update game status: %v", err)
		return
	}

	// Send match found message to both players
	matchMsg := map[string]interface{}{
		"type":       "match_found",
		"gameId":     game.ID,
		"seed":       game.Seed,
		"opponent":   game.Player2.Username,
		"opponentId": game.Player2.ID,
	}

	gm.sendToPlayer(game.Player1.ID, matchMsg)

	matchMsg["opponent"] = game.Player1.Username
	matchMsg["opponentId"] = game.Player1.ID
	gm.sendToPlayer(game.Player2.ID, matchMsg)

	log.Printf("Started game %s between %s and %s with seed %d",
		game.ID, game.Player1.Username, game.Player2.Username, game.Seed)
}

// HandleGameMove processes a player's move
func (gm *GameManager) HandleGameMove(playerID string, move *models.GameMove) error {
	// Get player
	player, err := gm.playerStore.GetPlayer(playerID)
	if err != nil {
		return err
	}

	// Get game
	game, err := gm.gameStore.GetGame(player.GameID)
	if err != nil {
		return err
	}

	// Validate player is in this game
	if game.Player1.ID != playerID && game.Player2.ID != playerID {
		return nil // Invalid player for this game
	}

	// Broadcast move to opponent
	var opponentID string
	if game.Player1.ID == playerID {
		opponentID = game.Player2.ID
	} else {
		opponentID = game.Player1.ID
	}

	moveMsg := map[string]interface{}{
		"type":      "game_move",
		"gameId":    game.ID,
		"playerId":  playerID,
		"moveType":  move.MoveType,
		"timestamp": move.Timestamp,
	}

	gm.sendToPlayer(opponentID, moveMsg)

	return nil
}

// HandleGameState processes a player's game state update
func (gm *GameManager) HandleGameState(playerID string, state *models.GameState) error {
	// Get player
	player, err := gm.playerStore.GetPlayer(playerID)
	if err != nil {
		return err
	}

	// Get game
	game, err := gm.gameStore.GetGame(player.GameID)
	if err != nil {
		return err
	}

	// Validate player is in this game
	if game.Player1.ID != playerID && game.Player2.ID != playerID {
		return nil // Invalid player for this game
	}

	// Broadcast state to opponent
	var opponentID string
	if game.Player1.ID == playerID {
		opponentID = game.Player2.ID
	} else {
		opponentID = game.Player1.ID
	}

	stateMsg := map[string]interface{}{
		"type":         "game_state",
		"gameId":       game.ID,
		"playerId":     playerID,
		"board":        state.Board,
		"score":        state.Score,
		"level":        state.Level,
		"lines":        state.Lines,
		"currentPiece": state.CurrentPiece,
		"timestamp":    state.Timestamp,
	}

	gm.sendToPlayer(opponentID, stateMsg)

	return nil
}

// EndGame handles game completion
func (gm *GameManager) EndGame(gameID, winnerID string) error {
	game, err := gm.gameStore.GetGame(gameID)
	if err != nil {
		return err
	}

	// Update game status
	game.Status = models.GameStatusFinished
	err = gm.gameStore.UpdateGame(game)
	if err != nil {
		return err
	}

	// Clear player game IDs
	game.Player1.GameID = ""
	game.Player2.GameID = ""
	gm.playerStore.UpdatePlayer(game.Player1)
	gm.playerStore.UpdatePlayer(game.Player2)

	// Send game over message
	gameOverMsg := map[string]interface{}{
		"type":     "game_over",
		"gameId":   gameID,
		"winnerId": winnerID,
	}

	gm.sendToPlayer(game.Player1.ID, gameOverMsg)
	gm.sendToPlayer(game.Player2.ID, gameOverMsg)

	log.Printf("Game %s ended, winner: %s", gameID, winnerID)

	return nil
}

// sendToPlayer sends a message to a specific player via WebSocket
func (gm *GameManager) sendToPlayer(playerID string, message map[string]interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	gm.wsManager.SendToPlayer(playerID, data)
}
