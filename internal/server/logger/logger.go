package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	// Create structured logger with JSON output
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// SetLevel sets the logging level
func SetLevel(level slog.Level) {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
