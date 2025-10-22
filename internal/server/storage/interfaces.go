package storage

import "github.com/briancain/go-tetris/pkg/models"

// PlayerStore handles player data persistence
type PlayerStore interface {
	CreatePlayer(player *models.Player) error
	GetPlayer(id string) (*models.Player, error)
	GetPlayerByToken(token string) (*models.Player, error)
	UpdatePlayer(player *models.Player) error
	DeletePlayer(id string) error
}

// GameStore handles game session persistence
type GameStore interface {
	CreateGame(game *models.GameSession) error
	GetGame(id string) (*models.GameSession, error)
	UpdateGame(game *models.GameSession) error
	DeleteGame(id string) error
	GetActiveGames() ([]*models.GameSession, error)
	GetAllGames() ([]*models.GameSession, error)
}

// QueueStore handles matchmaking queue
type QueueStore interface {
	AddToQueue(playerID string) error
	RemoveFromQueue(playerID string) error
	GetQueuedPlayers() ([]string, error)
	GetQueuePosition(playerID string) (int, error)
}
