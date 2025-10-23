package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/briancain/go-tetris/internal/server/services"
)

type HealthHandler struct {
	wsManager *services.WebSocketManager
}

type HealthResponse struct {
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Version     string            `json:"version"`
	Checks      map[string]string `json:"checks"`
	Connections int               `json:"websocket_connections"`
}

type MetricsResponse struct {
	WebSocketConnections int       `json:"websocket_connections"`
	Uptime               string    `json:"uptime"`
	Timestamp            time.Time `json:"timestamp"`
}

var startTime = time.Now()

func NewHealthHandler(wsManager *services.WebSocketManager) *HealthHandler {
	return &HealthHandler{
		wsManager: wsManager,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string)
	status := "healthy"

	// Check WebSocket manager
	if h.wsManager != nil {
		checks["websocket_manager"] = "ok"
	} else {
		checks["websocket_manager"] = "unavailable"
		status = "degraded"
	}

	// TODO: Add Redis connectivity check when Redis is implemented
	checks["redis"] = "not_implemented"

	response := HealthResponse{
		Status:      status,
		Timestamp:   time.Now(),
		Version:     "1.0.0",
		Checks:      checks,
		Connections: h.getConnectionCount(),
	}

	w.Header().Set("Content-Type", "application/json")
	if status != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)

	response := MetricsResponse{
		WebSocketConnections: h.getConnectionCount(),
		Uptime:               uptime.String(),
		Timestamp:            time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) getConnectionCount() int {
	if h.wsManager == nil {
		return 0
	}
	return h.wsManager.GetConnectionCount()
}
