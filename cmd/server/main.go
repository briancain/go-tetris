package main

import (
	"net/http"

	"github.com/briancain/go-tetris/internal/server/handlers"
	"github.com/briancain/go-tetris/internal/server/logger"
	"github.com/briancain/go-tetris/internal/server/middleware"
	"github.com/briancain/go-tetris/internal/server/services"
	"github.com/briancain/go-tetris/internal/server/storage/memory"
)

func main() {
	logger.Logger.Info("Starting Tetris multiplayer server")

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

	// Setup routes with logging middleware
	http.HandleFunc("/api/auth/login", middleware.RequestLogging(authHandler.Login))
	http.HandleFunc("/api/auth/logout", middleware.RequestLogging(authMiddleware.RequireAuth(authHandler.Logout)))

	http.HandleFunc("/api/matchmaking/queue", middleware.RequestLogging(authMiddleware.RequireAuth(matchmakingHandler.JoinQueue)))
	http.HandleFunc("/api/matchmaking/queue/leave", middleware.RequestLogging(authMiddleware.RequireAuth(matchmakingHandler.LeaveQueue)))
	http.HandleFunc("/api/matchmaking/status", middleware.RequestLogging(authMiddleware.RequireAuth(matchmakingHandler.GetQueueStatus)))

	http.HandleFunc("/api/leaderboard", middleware.RequestLogging(leaderboardHandler.GetLeaderboard))

	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Start server
	port := ":8080"
	logger.Logger.Info("Server starting",
		"port", port,
		"websocket_endpoint", "ws://localhost"+port+"/ws",
		"health_endpoint", "http://localhost"+port+"/health",
	)

	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Logger.Error("Server failed to start", "error", err)
	}
}
