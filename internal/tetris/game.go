package tetris

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Game states
const (
	StateMainMenu = iota
	StateMultiplayerSetup
	StateMatchmaking
	StatePlaying
	StatePaused
	StateGameOver
	StateRematchWaiting
	StateHighScores
)

// LeaderboardEntry represents a leaderboard entry
type LeaderboardEntry struct {
	Rank       int    `json:"rank"`
	Username   string `json:"username"`
	HighScore  int    `json:"highScore"`
	TotalGames int    `json:"totalGames"`
	Wins       int    `json:"wins"`
	Losses     int    `json:"losses"`
}

// Game represents the Tetris game state
type Game struct {
	Board             *Board
	CurrentPiece      *Piece
	NextPiece         *Piece
	HeldPiece         *Piece // Piece that is being held
	HasSwapped        bool   // Flag to prevent multiple swaps per turn
	State             int
	Score             int
	Level             int
	LinesCleared      int
	DropTimer         time.Time
	DropInterval      time.Duration
	LastMoveDown      time.Time
	LastMoveSide      time.Time
	LastRotate        time.Time
	LastHold          time.Time // Time of last hold action
	InputDelay        time.Duration
	FastDropDelay     time.Duration
	PieceGen          *PieceGenerator // 7-bag piece generator
	BackToBack        bool            // Track back-to-back special clears
	LastClearWasTSpin bool            // Track if the last clear was a T-spin
	LastWasBackToBack bool            // Track if the last clear got back-to-back bonus

	// Multiplayer fields
	MultiplayerMode   bool               `json:"multiplayerMode"`
	MultiplayerClient *MultiplayerClient `json:"-"` // Don't serialize WebSocket connection
	OpponentBoard     [][]Cell           `json:"opponentBoard,omitempty"`
	OpponentScore     int                `json:"opponentScore,omitempty"`
	OpponentLevel     int                `json:"opponentLevel,omitempty"`
	OpponentLines     int                `json:"opponentLines,omitempty"`
	LocalPlayerLost   bool               `json:"localPlayerLost,omitempty"`
	OpponentLost      bool               `json:"opponentLost,omitempty"`
	LoserScore        int                `json:"loserScore,omitempty"`
	RematchRequested  bool               `json:"rematchRequested,omitempty"`

	// UI state
	UsernameInput    string `json:"usernameInput,omitempty"`
	ConnectionStatus string `json:"connectionStatus,omitempty"`
	OpponentName     string `json:"opponentName,omitempty"`

	// Local high score (for single player)
	LocalHighScore int `json:"localHighScore,omitempty"`

	// Server leaderboard data
	Leaderboard []LeaderboardEntry `json:"leaderboard,omitempty"`
	ServerURL   string             `json:"serverURL,omitempty"`

	// Performance optimization: reusable slices
	boardBuffer     [][]Cell // Reusable board slice for multiplayer
	ghostY          int      // Cached ghost piece Y position
	ghostCacheValid bool     // Whether ghost cache is valid
}

// NewGame creates a new Tetris game
func NewGame() *Game {
	game := &Game{
		Board:             NewBoard(),
		State:             StateMainMenu,
		Score:             0,
		Level:             1,
		LinesCleared:      0,
		DropInterval:      800 * time.Millisecond, // Initial drop speed
		InputDelay:        100 * time.Millisecond, // Delay between input actions
		FastDropDelay:     50 * time.Millisecond,  // Fast drop speed
		PieceGen:          NewPieceGenerator(),    // Initialize the 7-bag generator
		BackToBack:        false,
		LastClearWasTSpin: false,
		ServerURL:         "http://localhost:8080", // Default server URL
	}

	// Initialize pieces
	game.NextPiece = game.PieceGen.NextPiece()

	return game
}

// NewGameWithSeed creates a new Tetris game with a specific random seed
func NewGameWithSeed(seed int64) *Game {
	game := &Game{
		Board:             NewBoard(),
		State:             StateMainMenu,
		Score:             0,
		Level:             1,
		LinesCleared:      0,
		DropInterval:      800 * time.Millisecond,          // Initial drop speed
		InputDelay:        100 * time.Millisecond,          // Delay between input actions
		FastDropDelay:     50 * time.Millisecond,           // Fast drop speed
		PieceGen:          NewPieceGeneratorWithSeed(seed), // Initialize with seed
		BackToBack:        false,
		LastClearWasTSpin: false,
	}

	// Initialize pieces
	game.NextPiece = game.PieceGen.NextPiece()

	return game
}

// Start begins a new game
func (g *Game) Start() {
	g.Board.Clear()
	g.HeldPiece = nil
	g.HasSwapped = false
	g.Score = 0
	g.Level = 1
	g.LinesCleared = 0
	g.State = StatePlaying
	g.DropTimer = time.Now()
	g.BackToBack = false
	g.LastClearWasTSpin = false
	g.LastWasBackToBack = false
	g.LocalPlayerLost = false
	g.OpponentLost = false
	g.LoserScore = 0
	g.RematchRequested = false
	g.updateDropInterval()

	// Generate fresh pieces and validate spawn position
	g.NextPiece = g.PieceGen.NextPiece()
	g.CurrentPiece = g.PieceGen.NextPiece()
	g.NextPiece = g.PieceGen.NextPiece()
	g.invalidateGhostCache()

	// Check for game over on initial spawn (important for rematch)
	if !g.Board.IsValidPosition(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y) {
		g.handleLocalGameOver()
	}
}

// canProcessInput returns true if the game can process input
func (g *Game) canProcessInput() bool {
	if g.State != StatePlaying {
		return false
	}
	// Stop input if local player lost in multiplayer
	if g.MultiplayerMode && g.LocalPlayerLost {
		return false
	}
	return true
}

// Update updates the game state
func (g *Game) Update() {
	// Process multiplayer messages in all states
	if g.MultiplayerMode {
		g.ProcessMultiplayerMessages()
	}

	if !g.canProcessInput() {
		return
	}

	// Check if it's time to drop the piece
	if time.Since(g.DropTimer) >= g.DropInterval {
		g.DropTimer = time.Now()
		g.moveDown()
	}
}

// invalidateGhostCache marks the ghost piece cache as invalid
func (g *Game) invalidateGhostCache() {
	g.ghostCacheValid = false
}

// moveDown moves the current piece down one row
func (g *Game) moveDown() {
	// Create a copy of the current piece and try to move it down
	testPiece := g.CurrentPiece.Copy()
	testPiece.Move(0, 1)

	// Check if the move is valid
	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Move(0, 1)
		return
	}

	// If we can't move down, lock the piece in place
	g.lockPiece()

	// Check for completed lines
	linesCleared := g.Board.ClearLines()
	if linesCleared > 0 {
		g.addScore(linesCleared)
	}

	// Spawn next piece and check for game over
	g.spawnNextPiece()
}

// MoveLeft moves the current piece left
func (g *Game) MoveLeft() bool {
	if !g.canProcessInput() || time.Since(g.LastMoveSide) < g.InputDelay {
		return false
	}

	g.LastMoveSide = time.Now()
	testPiece := g.CurrentPiece.Copy()
	testPiece.Move(-1, 0)

	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Move(-1, 0)
		g.invalidateGhostCache() // Invalidate ghost cache when piece moves
		g.sendMoveToServer("left")
		return true
	}

	return false
}

// MoveRight moves the current piece right
func (g *Game) MoveRight() bool {
	if !g.canProcessInput() || time.Since(g.LastMoveSide) < g.InputDelay {
		return false
	}

	g.LastMoveSide = time.Now()
	testPiece := g.CurrentPiece.Copy()
	testPiece.Move(1, 0)

	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Move(1, 0)
		g.invalidateGhostCache() // Invalidate ghost cache when piece moves
		g.sendMoveToServer("right")
		return true
	}

	return false
}

// RotatePiece rotates the current piece using SRS
func (g *Game) RotatePiece() bool {
	if !g.canProcessInput() || time.Since(g.LastRotate) < g.InputDelay {
		return false
	}

	// Stop input if local player lost in multiplayer
	if g.MultiplayerMode && g.LocalPlayerLost {
		return false
	}

	g.LastRotate = time.Now()

	// Skip rotation for O piece
	if g.CurrentPiece.Type == TypeO {
		return false
	}

	// Create a test piece for rotation
	testPiece := g.CurrentPiece.Copy()

	// Store the current rotation state
	oldRotationState := testPiece.RotationState

	// Rotate the test piece
	testPiece.Rotate()

	// Check if the basic rotation works
	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		// Apply the rotation to the actual piece
		g.CurrentPiece.Rotate()
		g.invalidateGhostCache() // Invalidate ghost cache when piece rotates
		g.sendMoveToServer("rotate")
		return true
	}

	// If basic rotation fails, try wall kicks
	var kickData [][]int

	// Get the appropriate kick data based on piece type and rotation transition
	if g.CurrentPiece.Type == TypeI {
		kickData = wallKickDataI[oldRotationState]
	} else {
		kickData = wallKickDataJLSTZ[oldRotationState]
	}

	// Try each wall kick
	for _, offset := range kickData {
		testX := testPiece.X + offset[0]
		testY := testPiece.Y + offset[1]

		if g.Board.IsValidPosition(testPiece, testX, testY) {
			// Apply the rotation and offset to the actual piece
			g.CurrentPiece.Rotate()
			g.CurrentPiece.X += offset[0]
			g.CurrentPiece.Y += offset[1]
			g.invalidateGhostCache() // Invalidate ghost cache when piece rotates
			g.sendMoveToServer("rotate")
			return true
		}
	}

	// If all wall kicks fail, don't rotate
	return false
}

// HardDrop drops the piece all the way down
func (g *Game) HardDrop() {
	if !g.canProcessInput() {
		return
	}

	// Get the ghost piece Y position
	ghostY := g.GetGhostPieceY()

	// Calculate how many cells we moved down
	distance := ghostY - g.CurrentPiece.Y

	// Move the piece to the ghost position
	g.CurrentPiece.Y = ghostY

	// Add score based on distance
	g.Score += distance

	g.sendMoveToServer("hard_drop")
	g.lockPiece()

	// Check for completed lines
	linesCleared := g.Board.ClearLines()
	if linesCleared > 0 {
		g.addScore(linesCleared)
	}

	// Send updated state to server
	g.sendStateToServer()

	// Spawn next piece and check for game over
	g.spawnNextPiece()
}

// SoftDrop accelerates the piece downward
func (g *Game) SoftDrop() bool {
	if !g.canProcessInput() || time.Since(g.LastMoveDown) < g.FastDropDelay {
		return false
	}

	g.LastMoveDown = time.Now()
	testPiece := g.CurrentPiece.Copy()
	testPiece.Move(0, 1)

	if g.Board.IsValidPosition(testPiece, testPiece.X, testPiece.Y) {
		g.CurrentPiece.Move(0, 1)
		g.Score++ // Small bonus for soft drop
		return true
	}

	return false
}

// TogglePause toggles the game's pause state
func (g *Game) TogglePause() {
	if g.State == StatePlaying {
		g.State = StatePaused
	} else if g.State == StatePaused {
		g.State = StatePlaying
		g.DropTimer = time.Now() // Reset drop timer when unpausing
	}
}

// lockPiece locks the current piece in place on the board
func (g *Game) lockPiece() {
	g.Board.PlacePiece(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y, true)
	g.HasSwapped = false // Reset swap flag when piece is locked

	// Check for game over - if any part of the piece locked in the hidden area (top 2 rows)
	shape := g.CurrentPiece.Shape
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[0]); j++ {
			if shape[i][j] {
				pieceY := g.CurrentPiece.Y + i
				// If any part of the piece is in the hidden area (rows 0-1), game over
				if pieceY < 2 {
					g.handleLocalGameOver()
					return
				}
			}
		}
	}
}

// handleLocalGameOver handles when the local player loses
func (g *Game) handleLocalGameOver() {
	// Update local high score for single player
	if !g.MultiplayerMode && g.Score > g.LocalHighScore {
		g.LocalHighScore = g.Score
		log.Printf("New local high score: %d", g.LocalHighScore)
	}

	if g.MultiplayerMode && g.MultiplayerClient != nil && g.MultiplayerClient.IsConnected() {
		// In multiplayer, notify server that we lost
		log.Printf("Game: Local player lost, notifying server")
		g.sendGameOverToServer()
		// Don't set StateGameOver yet - wait for server to end the match
	} else {
		// In single player, end immediately
		g.State = StateGameOver
	}
}

// sendGameOverToServer notifies the server that this player lost
func (g *Game) sendGameOverToServer() {
	if g.MultiplayerClient != nil && g.MultiplayerClient.IsConnected() {
		err := g.MultiplayerClient.SendGameOver()
		if err != nil {
			log.Printf("Failed to send game over message: %v", err)
		}
	}
}

// spawnNextPiece spawns the next piece and checks for game over
func (g *Game) spawnNextPiece() {
	g.CurrentPiece = g.NextPiece
	g.NextPiece = g.PieceGen.NextPiece()
	g.invalidateGhostCache() // Invalidate ghost cache for new piece

	// Check for game over - if the new piece can't be placed
	if !g.Board.IsValidPosition(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y) {
		g.handleLocalGameOver()
	}
}

// addScore adds to the score based on lines cleared and special moves
func (g *Game) addScore(linesCleared int) {
	// Check for T-spin
	isTSpin := g.isTSpin()

	// Base points for regular line clears
	baseLinePoints := []int{0, 100, 300, 500, 800}

	// Base points for T-spin line clears
	tSpinPoints := []int{400, 800, 1200, 1600}

	var points int

	// Special move detection
	isSpecialClear := false

	if isTSpin && linesCleared > 0 {
		// T-spin with line clear
		points = tSpinPoints[linesCleared-1] * g.Level
		isSpecialClear = true
		g.LastClearWasTSpin = true
	} else if linesCleared == 4 {
		// Tetris (4 lines)
		points = baseLinePoints[linesCleared] * g.Level
		isSpecialClear = true
		g.LastClearWasTSpin = false
	} else {
		// Regular line clear
		if linesCleared > 0 && linesCleared < len(baseLinePoints) {
			points = baseLinePoints[linesCleared] * g.Level
		}
		g.LastClearWasTSpin = false
	}

	// Back-to-Back bonus (50% bonus for consecutive special clears)
	g.LastWasBackToBack = false
	if isSpecialClear && g.BackToBack && linesCleared > 0 {
		points = points * 3 / 2 // 50% bonus
		g.LastWasBackToBack = true
	}

	// Update Back-to-Back status
	if linesCleared > 0 {
		g.BackToBack = isSpecialClear
	}

	g.Score += points
	g.LinesCleared += linesCleared

	// Level up every 10 lines
	newLevel := (g.LinesCleared / 10) + 1
	if newLevel > g.Level {
		g.Level = newLevel
		g.updateDropInterval()
	}
}

// updateDropInterval adjusts the piece drop speed based on level
func (g *Game) updateDropInterval() {
	// Formula: 800ms * (0.8^(level-1))
	// This makes the game get progressively faster with each level
	baseInterval := 800.0
	factor := 0.8

	// Calculate the new interval
	interval := baseInterval * pow(factor, float64(g.Level-1))
	g.DropInterval = time.Duration(interval) * time.Millisecond
}

// pow calculates x^y for our drop interval calculation
func pow(x, y float64) float64 {
	result := 1.0
	for i := 0; i < int(y); i++ {
		result *= x
	}
	return result
}

// HoldPiece swaps the current piece with the held piece
func (g *Game) HoldPiece() bool {
	if !g.canProcessInput() || time.Since(g.LastHold) < g.InputDelay {
		return false
	}

	// Can only swap once per piece
	if g.HasSwapped {
		return false
	}

	g.LastHold = time.Now()

	// If there's no held piece yet, store current piece and get next piece
	if g.HeldPiece == nil {
		g.HeldPiece = g.CurrentPiece.Copy()
		// Reset the held piece to its original position and orientation
		g.HeldPiece = NewPiece(g.HeldPiece.Type)

		g.CurrentPiece = g.NextPiece
		g.NextPiece = g.PieceGen.NextPiece()
	} else {
		// Swap current piece with held piece
		tempPiece := g.CurrentPiece

		// Create a new piece of the held type at the top of the board
		g.CurrentPiece = NewPiece(g.HeldPiece.Type)

		// Store the previous current piece as held
		g.HeldPiece = NewPiece(tempPiece.Type)
	}

	// Mark that we've swapped this turn
	g.HasSwapped = true

	// Check if the new current piece can be placed
	if !g.Board.IsValidPosition(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y) {
		// Game over if the piece can't be placed
		g.State = StateGameOver
		return false
	}

	return true
}

// isTSpin determines if the last move was a T-spin
// Uses the 3-corner T rule: A T-spin occurs when the T piece is rotated
// and at least 3 of the 4 corners surrounding the center of the T are occupied
func (g *Game) isTSpin() bool {
	// Only T pieces can perform T-spins
	if g.CurrentPiece.Type != TypeT {
		return false
	}

	// Get the center position of the T piece
	centerX := g.CurrentPiece.X + 1
	centerY := g.CurrentPiece.Y + 1

	// Check the four corners around the center
	cornerCount := 0

	// Top-left corner
	if centerX-1 < 0 || centerY-1 < 0 ||
		centerY-1 >= BoardHeightWithBuffer ||
		(centerX-1 < BoardWidth && g.Board.Cells[centerY-1][centerX-1] != Empty) {
		cornerCount++
	}

	// Top-right corner
	if centerX+1 >= BoardWidth || centerY-1 < 0 ||
		centerY-1 >= BoardHeightWithBuffer ||
		(centerX+1 < BoardWidth && centerY-1 < BoardHeightWithBuffer && g.Board.Cells[centerY-1][centerX+1] != Empty) {
		cornerCount++
	}

	// Bottom-left corner
	if centerX-1 < 0 || centerY+1 >= BoardHeightWithBuffer ||
		(centerX-1 < BoardWidth && centerY+1 < BoardHeightWithBuffer && g.Board.Cells[centerY+1][centerX-1] != Empty) {
		cornerCount++
	}

	// Bottom-right corner
	if centerX+1 >= BoardWidth || centerY+1 >= BoardHeightWithBuffer ||
		(centerX+1 < BoardWidth && centerY+1 < BoardHeightWithBuffer && g.Board.Cells[centerY+1][centerX+1] != Empty) {
		cornerCount++
	}

	// T-spin requires at least 3 corners to be occupied
	return cornerCount >= 3
}

// GetGhostPieceY calculates where the current piece would land if dropped
func (g *Game) GetGhostPieceY() int {
	if g.CurrentPiece == nil {
		return 0
	}

	// Return cached value if still valid
	if g.ghostCacheValid {
		return g.ghostY
	}

	// Calculate ghost position without creating piece copy
	ghostY := g.CurrentPiece.Y

	for {
		testY := ghostY + 1

		// Check if position is valid without creating a copy
		if !g.Board.IsValidPosition(g.CurrentPiece, g.CurrentPiece.X, testY) {
			break
		}

		ghostY = testY
	}

	// Cache the result
	g.ghostY = ghostY
	g.ghostCacheValid = true

	return ghostY
}

// clearOpponentBoard clears the opponent board without reallocating
func (g *Game) clearOpponentBoard() {
	for i := range g.OpponentBoard {
		for j := range g.OpponentBoard[i] {
			g.OpponentBoard[i][j] = Empty
		}
	}
}

// EnableMultiplayer enables multiplayer mode with server connection
func (g *Game) EnableMultiplayer(serverURL string) error {
	if g.MultiplayerClient != nil {
		g.MultiplayerClient.Close()
	}

	g.MultiplayerClient = NewMultiplayerClient(serverURL)
	g.MultiplayerMode = true

	// Initialize opponent board only if not already allocated
	if g.OpponentBoard == nil {
		g.OpponentBoard = make([][]Cell, BoardHeightWithBuffer)
		for i := range g.OpponentBoard {
			g.OpponentBoard[i] = make([]Cell, BoardWidth)
		}
	} else {
		// Reuse existing board, just clear it
		g.clearOpponentBoard()
	}

	return nil
}

// ConnectToServer connects to the multiplayer server
func (g *Game) ConnectToServer(username string) error {
	if g.MultiplayerClient == nil {
		return fmt.Errorf("multiplayer not enabled")
	}

	err := g.MultiplayerClient.Login(username)
	if err != nil {
		return err
	}

	return g.MultiplayerClient.Connect()
}

// JoinMatchmaking joins the matchmaking queue
func (g *Game) JoinMatchmaking() error {
	if g.MultiplayerClient == nil {
		return fmt.Errorf("multiplayer not enabled")
	}

	return g.MultiplayerClient.JoinQueue()
}

// ProcessMultiplayerMessages processes incoming server messages
func (g *Game) ProcessMultiplayerMessages() {
	if g.MultiplayerClient == nil {
		return
	}

	for {
		message := g.MultiplayerClient.GetMessage()
		if message == nil {
			break // No more messages
		}

		g.handleMultiplayerMessage(message)
	}
}

// HandleMultiplayerMessage handles a single multiplayer message (public for testing)
func (g *Game) HandleMultiplayerMessage(message map[string]interface{}) {
	g.handleMultiplayerMessage(message)
}

// handleMultiplayerMessage handles a single multiplayer message (internal)
func (g *Game) handleMultiplayerMessage(message map[string]interface{}) {
	msgType, ok := message["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "match_found":
		g.handleMatchFound(message)
	case "game_move":
		g.handleOpponentMove(message)
	case "game_state":
		g.handleOpponentState(message)
	case "game_over":
		g.handleGameOver(message)
	case "player_lost":
		g.handlePlayerLost(message)
	case "rematch_request":
		g.handleRematchRequest(message)
	case "rematch_start":
		g.handleRematchStart(message)
	case "opponent_disconnected":
		g.handleOpponentDisconnected(message)
	}
}

// handleMatchFound processes match found message
func (g *Game) handleMatchFound(message map[string]interface{}) {
	seed, ok := message["seed"].(float64)
	if ok {
		// Use server-provided seed
		g.PieceGen = NewPieceGeneratorWithSeed(int64(seed))
		g.NextPiece = g.PieceGen.NextPiece()
		log.Printf("Game: Using server seed: %.0f", seed)
	}

	opponent, ok := message["opponent"].(string)
	if ok {
		g.OpponentName = opponent
		log.Printf("Game: Matched with opponent: %s", opponent)
	}

	// Start the game - this will change state to StatePlaying
	g.Start()
}

// handleOpponentMove processes opponent move
func (g *Game) handleOpponentMove(message map[string]interface{}) {
	moveType, ok := message["moveType"].(string)
	if ok {
		log.Printf("Game: Opponent move: %s", moveType)
		// In a full implementation, you might want to show opponent moves visually
	}
}

// handleOpponentState processes opponent game state
func (g *Game) handleOpponentState(message map[string]interface{}) {
	// Update opponent score
	if score, ok := message["score"].(float64); ok {
		g.OpponentScore = int(score)
	}

	// Update opponent level
	if level, ok := message["level"].(float64); ok {
		g.OpponentLevel = int(level)
	}

	// Update opponent lines
	if lines, ok := message["lines"].(float64); ok {
		g.OpponentLines = int(lines)
	}

	// Update opponent board
	if boardInterface, ok := message["board"].([]interface{}); ok {
		for i, rowInterface := range boardInterface {
			if i >= len(g.OpponentBoard) {
				break
			}
			if row, ok := rowInterface.([]interface{}); ok {
				for j, cellInterface := range row {
					if j >= len(g.OpponentBoard[i]) {
						break
					}
					if cellValue, ok := cellInterface.(float64); ok {
						g.OpponentBoard[i][j] = Cell(int(cellValue))
					}
				}
			}
		}
	}
}

// handleGameOver processes game over message
func (g *Game) handleGameOver(message map[string]interface{}) {
	winnerID, ok := message["winnerId"].(string)
	if ok {
		log.Printf("Game: Game over, winner: %s", winnerID)
	}

	// End the game
	g.State = StateGameOver
}

// handlePlayerLost processes when a player loses but game continues
func (g *Game) handlePlayerLost(message map[string]interface{}) {
	playerID, ok := message["playerId"].(string)
	if !ok {
		return
	}

	// Get loser's score
	if loserScore, ok := message["loserScore"].(float64); ok {
		g.LoserScore = int(loserScore)
	}

	if g.MultiplayerClient != nil && playerID == g.MultiplayerClient.playerID {
		// We lost - enter spectator mode but don't end game yet
		g.LocalPlayerLost = true
		log.Printf("Game: Local player lost (score: %d), waiting for opponent to beat score", g.LoserScore)
	} else {
		// Opponent lost - we continue playing until we beat their score
		g.OpponentLost = true
		log.Printf("Game: Opponent lost (score: %d), continue playing to beat their score", g.LoserScore)
	}
}

// sendMoveToServer sends a move to the server
func (g *Game) sendMoveToServer(moveType string) {
	if g.MultiplayerClient != nil && g.MultiplayerClient.IsConnected() {
		_ = g.MultiplayerClient.SendGameMove(moveType)
	}
}

// RequestRematch sends a rematch request to the server
func (g *Game) RequestRematch() {
	if !g.MultiplayerMode || g.MultiplayerClient == nil {
		return
	}

	g.RematchRequested = true
	g.State = StateRematchWaiting

	message := map[string]interface{}{
		"type": "rematch_request",
	}

	if g.MultiplayerClient.IsConnected() {
		_ = g.MultiplayerClient.sendMessage(message)
	}

	log.Printf("Game: Rematch requested")
}

// handleRematchRequest processes rematch request from opponent
func (g *Game) handleRematchRequest(_ map[string]interface{}) {
	log.Printf("Game: Opponent requested rematch")
	// Could show UI notification here
}

// handleRematchStart processes rematch start from server
func (g *Game) handleRematchStart(message map[string]interface{}) {
	seed, ok := message["seed"].(float64)
	if ok {
		g.PieceGen.SetSeed(int64(seed))
		log.Printf("Game: Rematch starting with seed: %d", int64(seed))
	}

	// Reset game state for rematch
	g.Start()
	log.Printf("Game: Rematch started")
}

// FetchLeaderboard fetches the server leaderboard
func (g *Game) FetchLeaderboard() {
	// Make HTTP request to leaderboard endpoint
	resp, err := http.Get(g.ServerURL + "/api/leaderboard?limit=10")
	if err != nil {
		log.Printf("Failed to fetch leaderboard (server may not be running): %v", err)
		// Set empty leaderboard to show "No server scores available"
		g.Leaderboard = []LeaderboardEntry{}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Leaderboard request failed with status: %d", resp.StatusCode)
		g.Leaderboard = []LeaderboardEntry{}
		return
	}

	var leaderboard []LeaderboardEntry
	err = json.NewDecoder(resp.Body).Decode(&leaderboard)
	if err != nil {
		log.Printf("Failed to parse leaderboard response: %v", err)
		g.Leaderboard = []LeaderboardEntry{}
		return
	}

	g.Leaderboard = leaderboard
	log.Printf("Fetched leaderboard with %d entries", len(leaderboard))
}

// handleOpponentDisconnected processes opponent disconnect message
func (g *Game) handleOpponentDisconnected(_ map[string]interface{}) {
	log.Printf("Game: Opponent disconnected - You win!")
	g.State = StateGameOver
}

// sendStateToServer sends current game state to server
func (g *Game) sendStateToServer() {
	if g.MultiplayerClient != nil && g.MultiplayerClient.IsConnected() {
		// Reuse board buffer to avoid allocations
		if g.boardBuffer == nil {
			g.boardBuffer = make([][]Cell, len(g.Board.Cells))
			for i := range g.boardBuffer {
				g.boardBuffer[i] = make([]Cell, len(g.Board.Cells[i]))
			}
		}

		// Copy current board state to buffer
		for i, row := range g.Board.Cells {
			copy(g.boardBuffer[i], row[:])
		}

		_ = g.MultiplayerClient.SendGameState(g.boardBuffer, g.Score, g.Level, g.LinesCleared)
	}
}

// AddToUsernameInput adds a character to the username input
func (g *Game) AddToUsernameInput(char rune) {
	if len(g.UsernameInput) < 12 { // Limit to 12 characters for better display
		g.UsernameInput += string(char)
	}
}

// RemoveFromUsernameInput removes the last character from username input
func (g *Game) RemoveFromUsernameInput() {
	if len(g.UsernameInput) > 0 {
		g.UsernameInput = g.UsernameInput[:len(g.UsernameInput)-1]
	}
}

// StartMultiplayerConnection attempts to connect to the multiplayer server
func (g *Game) StartMultiplayerConnection() {
	if len(g.UsernameInput) < 2 {
		g.ConnectionStatus = "Username must be at least 2 characters"
		return
	}
	if len(g.UsernameInput) > 12 {
		g.ConnectionStatus = "Username too long (max 12 characters)"
		return
	}

	g.ConnectionStatus = "Connecting to server..."

	// Enable multiplayer mode
	err := g.EnableMultiplayer("http://localhost:8080")
	if err != nil {
		log.Printf("Multiplayer: Failed to enable multiplayer: %v", err)
		g.ConnectionStatus = "Connection error occurred"
		return
	}

	// Connect to server
	err = g.ConnectToServer(g.UsernameInput)
	if err != nil {
		log.Printf("Multiplayer: Failed to connect to server: %v", err)
		g.ConnectionStatus = err.Error()
		return
	}

	// Join matchmaking
	err = g.JoinMatchmaking()
	if err != nil {
		log.Printf("Multiplayer: Failed to join matchmaking queue: %v", err)
		g.ConnectionStatus = "Connection error occurred"
		return
	}

	// Success - move to matchmaking state
	log.Printf("Multiplayer: Successfully connected as %s", g.UsernameInput)
	g.ConnectionStatus = "Finding match..."
	g.State = StateMatchmaking
}
