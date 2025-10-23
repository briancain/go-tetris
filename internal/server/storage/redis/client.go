package redis

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/briancain/go-tetris/internal/server/logger"
)

// redisLogger adapts Redis client logging to our structured logger
type redisLogger struct{}

func (l *redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	msg := strings.TrimSpace(format)

	// Convert Redis log levels to our structured logging
	if strings.Contains(msg, "error") || strings.Contains(msg, "ERROR") {
		logger.Logger.Error("Redis client", "message", msg, "args", v)
	} else if strings.Contains(msg, "warn") || strings.Contains(msg, "WARN") {
		logger.Logger.Warn("Redis client", "message", msg, "args", v)
	} else {
		logger.Logger.Debug("Redis client", "message", msg, "args", v)
	}
}

// Client wraps Redis client with health check capability
type Client struct {
	*redis.Client
}

// NewClient creates a new Redis client with connection pooling
func NewClient(redisURL string) (*Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	// Configure connection pooling
	opts.PoolSize = 10
	opts.MinIdleConns = 2
	opts.MaxIdleConns = 5
	opts.ConnMaxIdleTime = 5 * time.Minute

	client := redis.NewClient(opts)

	// Use our structured logger for Redis client logs
	redis.SetLogger(&redisLogger{})

	return &Client{Client: client}, nil
}

// HealthCheck implements storage.HealthChecker
func (c *Client) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return c.Ping(ctx).Err()
}
