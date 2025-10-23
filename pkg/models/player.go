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
	// Stats
	TotalGames int `json:"totalGames"`
	Wins       int `json:"wins"`
	Losses     int `json:"losses"`
	HighScore  int `json:"highScore"`
}
