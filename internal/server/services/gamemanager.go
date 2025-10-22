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

	// Update player's score in game session
	if game.Player1.ID == playerID {
		game.Player1Score = state.Score
	} else {
		game.Player2Score = state.Score
	}

	// Check if surviving player has won by score
	gm.checkScoreWin(game)

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

// EndGame handles when a player loses
func (gm *GameManager) EndGame(gameID, loserID string) error {
	game, err := gm.gameStore.GetGame(gameID)
	if err != nil {
		return err
	}

	// Mark player as lost and get their final score
	var loserScore int
	if game.Player1.ID == loserID {
		game.Player1Lost = true
		loserScore = game.Player1Score
	} else if game.Player2.ID == loserID {
		game.Player2Lost = true
		loserScore = game.Player2Score
	}

	// Update game in storage
	err = gm.gameStore.UpdateGame(game)
	if err != nil {
		return err
	}

	// Send player lost message
	playerLostMsg := map[string]interface{}{
		"type":       "player_lost",
		"gameId":     gameID,
		"playerId":   loserID,
		"loserScore": loserScore,
	}

	gm.sendToPlayer(game.Player1.ID, playerLostMsg)
	gm.sendToPlayer(game.Player2.ID, playerLostMsg)

	// Check if both players lost
	if game.Player1Lost && game.Player2Lost {
		gm.finalizeGame(game, "")
	}

	log.Printf("Player %s lost in game %s (score: %d)", loserID, gameID, loserScore)
	return nil
}

// checkScoreWin checks if surviving player has won by achieving higher score
func (gm *GameManager) checkScoreWin(game *models.GameSession) {
	// Only check if exactly one player has lost
	if game.Player1Lost && !game.Player2Lost {
		// Player 1 lost, check if Player 2 has beaten their score
		if game.Player2Score > game.Player1Score {
			gm.finalizeGame(game, game.Player2.ID)
		}
	} else if game.Player2Lost && !game.Player1Lost {
		// Player 2 lost, check if Player 1 has beaten their score
		if game.Player1Score > game.Player2Score {
			gm.finalizeGame(game, game.Player1.ID)
		}
	}
}

// finalizeGame ends the game with final results
func (gm *GameManager) finalizeGame(game *models.GameSession, winnerID string) {
	// Update game status
	game.Status = models.GameStatusFinished
	err := gm.gameStore.UpdateGame(game)
	if err != nil {
		log.Printf("Failed to update game status: %v", err)
		return
	}

	// Clear player game IDs
	game.Player1.GameID = ""
	game.Player2.GameID = ""
	gm.playerStore.UpdatePlayer(game.Player1)
	gm.playerStore.UpdatePlayer(game.Player2)

	// Send final game over message
	gameOverMsg := map[string]interface{}{
		"type":         "game_over",
		"gameId":       game.ID,
		"winnerId":     winnerID,
		"final":        true,
		"player1Score": game.Player1Score,
		"player2Score": game.Player2Score,
	}

	gm.sendToPlayer(game.Player1.ID, gameOverMsg)
	gm.sendToPlayer(game.Player2.ID, gameOverMsg)

	log.Printf("Game %s finalized, winner: %s (P1: %d, P2: %d)",
		game.ID, winnerID, game.Player1Score, game.Player2Score)
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
