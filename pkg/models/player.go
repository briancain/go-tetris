package models

import "time"

// Player represents a connected player
type Player struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	SessionToken string    `json:"sessionToken"`
	ConnectedAt  time.Time `json:"connectedAt"`
	InQueue      bool      `json:"inQueue"`
	GameID       string    `json:"gameId,omitempty"`
}

// GameSession represents an active game between two players
type GameSession struct {
	ID        string     `json:"id"`
	Player1   *Player    `json:"player1"`
	Player2   *Player    `json:"player2"`
	Seed      int64      `json:"seed"`
	Status    GameStatus `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
}

// GameStatus represents the current state of a game
type GameStatus string

const (
	GameStatusWaiting  GameStatus = "waiting"
	GameStatusActive   GameStatus = "active"
	GameStatusFinished GameStatus = "finished"
)

// GameMove represents a player's move
type GameMove struct {
	PlayerID  string    `json:"playerId"`
	GameID    string    `json:"gameId"`
	MoveType  string    `json:"moveType"` // "left", "right", "rotate", "drop", "hold"
	Timestamp time.Time `json:"timestamp"`
}

// GameState represents the current state of a player's board
type GameState struct {
	PlayerID     string    `json:"playerId"`
	GameID       string    `json:"gameId"`
	Board        [][]int   `json:"board"`
	Score        int       `json:"score"`
	Level        int       `json:"level"`
	Lines        int       `json:"lines"`
	CurrentPiece *Piece    `json:"currentPiece,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// Piece represents a Tetris piece for state sync
type Piece struct {
	Type     int      `json:"type"`
	X        int      `json:"x"`
	Y        int      `json:"y"`
	Rotation int      `json:"rotation"`
	Shape    [][]bool `json:"shape"`
}
