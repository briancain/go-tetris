package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/briancain/go-tetris/internal/server/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestLogging middleware logs HTTP requests and responses
func RequestLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID
		requestID := generateRequestID()

		// Add request ID to context
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		r = r.WithContext(ctx)

		// Wrap response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Log request
		logger.Logger.Info("HTTP request started",
			"requestID", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remoteAddr", r.RemoteAddr,
			"userAgent", r.UserAgent(),
		)

		// Call next handler
		next(rw, r)

		// Log response
		duration := time.Since(start)
		logger.Logger.Info("HTTP request completed",
			"requestID", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"statusCode", rw.statusCode,
			"duration", duration.String(),
		)
	}
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
