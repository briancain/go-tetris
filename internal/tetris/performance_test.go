package tetris

import (
	"testing"

	_ "github.com/briancain/go-tetris/internal/testutil" // Import for init side effects
)

func TestGhostCacheOptimization(t *testing.T) {
	game := NewGame()
	game.Start()

	// First call should calculate and cache
	ghostY1 := game.GetGhostPieceY()
	if !game.ghostCacheValid {
		t.Error("Ghost cache should be valid after first calculation")
	}

	// Second call should use cache (same result)
	ghostY2 := game.GetGhostPieceY()
	if ghostY1 != ghostY2 {
		t.Errorf("Ghost Y should be same from cache: %d != %d", ghostY1, ghostY2)
	}

	// Moving piece should invalidate cache
	game.MoveLeft()
	if game.ghostCacheValid {
		t.Error("Ghost cache should be invalid after piece moves")
	}

	// Rotating piece should invalidate cache
	game.ghostCacheValid = true // Manually set to test rotation
	game.RotatePiece()
	if game.ghostCacheValid {
		t.Error("Ghost cache should be invalid after piece rotates")
	}
}

func TestOpponentBoardReuse(t *testing.T) {
	game := NewGame()

	// First enable should allocate
	err := game.EnableMultiplayer("http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to enable multiplayer: %v", err)
	}

	if game.OpponentBoard == nil {
		t.Error("Opponent board should be allocated")
	}

	// Store pointer to first slice for comparison
	originalBoardPtr := &game.OpponentBoard[0][0]

	// Second enable should reuse, not reallocate
	err = game.EnableMultiplayer("http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to enable multiplayer second time: %v", err)
	}

	// Check if same memory is reused by comparing first element address
	newBoardPtr := &game.OpponentBoard[0][0]
	if originalBoardPtr != newBoardPtr {
		t.Error("Opponent board should be reused, not reallocated")
	}

	// Board should be cleared
	for i := range game.OpponentBoard {
		for j := range game.OpponentBoard[i] {
			if game.OpponentBoard[i][j] != Empty {
				t.Errorf("Opponent board should be cleared at [%d][%d]", i, j)
			}
		}
	}
}

func TestBoardBufferReuse(t *testing.T) {
	game := NewGame()
	game.Start()

	// Enable multiplayer to initialize board buffer
	err := game.EnableMultiplayer("http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to enable multiplayer: %v", err)
	}

	// Manually allocate the board buffer to test reuse
	game.boardBuffer = make([][]Cell, len(game.Board.Cells))
	for i := range game.boardBuffer {
		game.boardBuffer[i] = make([]Cell, len(game.Board.Cells[i]))
	}

	if game.boardBuffer == nil {
		t.Error("Board buffer should be allocated")
	}

	// Store pointer to first element for comparison
	if len(game.boardBuffer) > 0 && len(game.boardBuffer[0]) > 0 {
		originalBufferPtr := &game.boardBuffer[0][0]

		// Call sendStateToServer which should reuse the existing buffer
		game.sendStateToServer()

		newBufferPtr := &game.boardBuffer[0][0]
		if originalBufferPtr != newBufferPtr {
			t.Error("Board buffer should be reused, not reallocated")
		}
	}
}
