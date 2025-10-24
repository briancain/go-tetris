package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/briancain/go-tetris/internal/server/config"
	"github.com/briancain/go-tetris/internal/server/handlers"
	"github.com/briancain/go-tetris/internal/server/logger"
	"github.com/briancain/go-tetris/internal/server/middleware"
	"github.com/briancain/go-tetris/internal/server/services"
	"github.com/briancain/go-tetris/internal/server/storage"
	"github.com/briancain/go-tetris/internal/server/storage/memory"
	"github.com/briancain/go-tetris/internal/server/storage/redis"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger.Logger.Info("Starting Tetris multiplayer server", "config", cfg)

	// Initialize storage based on configuration
	var gameStore storage.GameStore
	var queueStore storage.QueueStore
	var storageHealth storage.HealthChecker

	playerStore := memory.NewPlayerStore()

	if cfg.RedisURL != "" {
		// Use Redis for game and queue storage
		logger.Logger.Info("Storage mode: Redis", "redis_url", cfg.RedisURL, "components", "games,queues")

		redisClient, err := redis.NewClient(cfg.RedisURL)
		if err != nil {
			logger.Logger.Error("Failed to connect to Redis", "error", err)
			os.Exit(1)
		}

		// Test Redis connection
		if err := redisClient.HealthCheck(); err != nil {
			logger.Logger.Error("Redis health check failed", "error", err)
			os.Exit(1)
		}

		gameStore = redis.NewGameStore(redisClient)
		queueStore = redis.NewQueueStore(redisClient)
		storageHealth = redisClient

		logger.Logger.Info("Redis storage initialized successfully")
	} else {
		// Use in-memory storage
		logger.Logger.Info("Storage mode: In-Memory", "components", "games,queues,players")
		gameStore = memory.NewGameStore()
		queueStore = memory.NewQueueStore()
		storageHealth = playerStore
	}

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
	healthHandler := handlers.NewHealthHandler(wsManager, storageHealth)

	// Setup routes with logging middleware
	http.HandleFunc("/api/auth/login", middleware.RequestLogging(authHandler.Login))
	http.HandleFunc("/api/auth/logout", middleware.RequestLogging(authMiddleware.RequireAuth(authHandler.Logout)))

	http.HandleFunc("/api/matchmaking/queue", middleware.RequestLogging(authMiddleware.RequireAuth(matchmakingHandler.JoinQueue)))
	http.HandleFunc("/api/matchmaking/queue/leave", middleware.RequestLogging(authMiddleware.RequireAuth(matchmakingHandler.LeaveQueue)))
	http.HandleFunc("/api/matchmaking/status", middleware.RequestLogging(authMiddleware.RequireAuth(matchmakingHandler.GetQueueStatus)))

	http.HandleFunc("/api/leaderboard", middleware.RequestLogging(leaderboardHandler.GetLeaderboard))

	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Health check and metrics endpoints
	http.HandleFunc("/health", healthHandler.Health)
	http.HandleFunc("/metrics", healthHandler.Metrics)

	// Create HTTP server
	port := ":" + cfg.Port
	server := &http.Server{
		Addr: port,
	}

	// Start server in goroutine
	go func() {
		logger.Logger.Info("Server starting",
			"port", port,
			"redis_url", cfg.RedisURL,
			"server_url", cfg.ServerURL,
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-quit
	logger.Logger.Info("Shutdown signal received, starting graceful shutdown...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown WebSocket connections
	logger.Logger.Info("Closing WebSocket connections...")
	wsManager.Shutdown()

	// Shutdown HTTP server
	logger.Logger.Info("Shutting down HTTP server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Logger.Info("Server gracefully stopped")
}
