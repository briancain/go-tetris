package main

import (
	"log"
	"net/http"

	"github.com/briancain/go-tetris/internal/server/handlers"
	"github.com/briancain/go-tetris/internal/server/middleware"
	"github.com/briancain/go-tetris/internal/server/services"
	"github.com/briancain/go-tetris/internal/server/storage/memory"
)

func main() {
	log.Println("Starting Tetris multiplayer server...")

	// Initialize storage
	playerStore := memory.NewPlayerStore()
	gameStore := memory.NewGameStore()
	queueStore := memory.NewQueueStore()

	// Initialize services
	authService := services.NewAuthService(playerStore)
	wsManager := services.NewWebSocketManager()
	gameManager := services.NewGameManager(gameStore, playerStore, wsManager)
	matchmakingService := services.NewMatchmakingService(playerStore, gameStore, queueStore, gameManager)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	matchmakingHandler := handlers.NewMatchmakingHandler(matchmakingService)
	leaderboardHandler := handlers.NewLeaderboardHandler(playerStore)
	wsHandler := handlers.NewWebSocketHandler(wsManager, authService, gameManager)

	// Setup routes
	http.HandleFunc("/api/auth/login", authHandler.Login)
	http.HandleFunc("/api/auth/logout", authMiddleware.RequireAuth(authHandler.Logout))

	http.HandleFunc("/api/matchmaking/queue", authMiddleware.RequireAuth(matchmakingHandler.JoinQueue))
	http.HandleFunc("/api/matchmaking/queue/leave", authMiddleware.RequireAuth(matchmakingHandler.LeaveQueue))
	http.HandleFunc("/api/matchmaking/status", authMiddleware.RequireAuth(matchmakingHandler.GetQueueStatus))

	http.HandleFunc("/api/leaderboard", leaderboardHandler.GetLeaderboard)

	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("WebSocket endpoint: ws://localhost%s/ws", port)
	log.Printf("Health check: http://localhost%s/health", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
