package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/briancain/go-tetris/internal/server/services"
)

// mockHealthChecker for testing
type mockHealthChecker struct {
	shouldFail bool
}

func (m *mockHealthChecker) HealthCheck() error {
	if m.shouldFail {
		return errors.New("storage connection failed")
	}
	return nil
}

func TestHealthHandler_Health(t *testing.T) {
	wsManager := services.NewWebSocketManager()
	healthChecker := &mockHealthChecker{}
	handler := NewHealthHandler(wsManager, healthChecker)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}

	if response.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", response.Version)
	}

	if response.Connections != 0 {
		t.Errorf("Expected 0 connections, got %d", response.Connections)
	}

	if response.Checks["websocket_manager"] != "ok" {
		t.Errorf("Expected websocket_manager check to be 'ok', got '%s'", response.Checks["websocket_manager"])
	}

	if response.Checks["storage"] != "ok" {
		t.Errorf("Expected storage check to be 'ok', got '%s'", response.Checks["storage"])
	}
}

func TestHealthHandler_HealthWithNilWebSocketManager(t *testing.T) {
	healthChecker := &mockHealthChecker{}
	handler := NewHealthHandler(nil, healthChecker)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", w.Code)
	}

	var response HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "degraded" {
		t.Errorf("Expected status 'degraded', got '%s'", response.Status)
	}

	if response.Connections != 0 {
		t.Errorf("Expected 0 connections, got %d", response.Connections)
	}

	if response.Checks["websocket_manager"] != "unavailable" {
		t.Errorf("Expected websocket_manager check to be 'unavailable', got '%s'", response.Checks["websocket_manager"])
	}
}

func TestHealthHandler_Metrics(t *testing.T) {
	wsManager := services.NewWebSocketManager()
	healthChecker := &mockHealthChecker{}
	handler := NewHealthHandler(wsManager, healthChecker)

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	handler.Metrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response MetricsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.WebSocketConnections != 0 {
		t.Errorf("Expected 0 WebSocket connections, got %d", response.WebSocketConnections)
	}

	if response.Uptime == "" {
		t.Error("Expected uptime to be set")
	}

	if response.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

func TestHealthHandler_MetricsWithNilWebSocketManager(t *testing.T) {
	healthChecker := &mockHealthChecker{}
	handler := NewHealthHandler(nil, healthChecker)

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	handler.Metrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response MetricsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.WebSocketConnections != 0 {
		t.Errorf("Expected 0 WebSocket connections, got %d", response.WebSocketConnections)
	}
}

func TestHealthHandler_ContentType(t *testing.T) {
	wsManager := services.NewWebSocketManager()
	healthChecker := &mockHealthChecker{}
	handler := NewHealthHandler(wsManager, healthChecker)

	tests := []struct {
		name     string
		endpoint string
		handler  http.HandlerFunc
	}{
		{"health", "/health", handler.Health},
		{"metrics", "/metrics", handler.Metrics},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.endpoint, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}
		})
	}
}

func TestHealthHandler_FailureScenarios(t *testing.T) {
	tests := []struct {
		name           string
		wsManager      *services.WebSocketManager
		expectedStatus int
		expectedHealth string
	}{
		{
			name:           "healthy with websocket manager",
			wsManager:      services.NewWebSocketManager(),
			expectedStatus: http.StatusOK,
			expectedHealth: "healthy",
		},
		{
			name:           "degraded without websocket manager",
			wsManager:      nil,
			expectedStatus: http.StatusServiceUnavailable,
			expectedHealth: "degraded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			healthChecker := &mockHealthChecker{}
			handler := NewHealthHandler(tt.wsManager, healthChecker)
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			handler.Health(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response HealthResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Status != tt.expectedHealth {
				t.Errorf("Expected health status '%s', got '%s'", tt.expectedHealth, response.Status)
			}
		})
	}
}

func TestHealthHandler_StorageHealthCheckFailure(t *testing.T) {
	wsManager := services.NewWebSocketManager()
	healthChecker := &mockHealthChecker{shouldFail: true}
	handler := NewHealthHandler(wsManager, healthChecker)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var response HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "unhealthy" {
		t.Errorf("Expected status 'unhealthy', got '%s'", response.Status)
	}

	if response.Checks["storage"] != "error: storage connection failed" {
		t.Errorf("Expected storage error message, got '%s'", response.Checks["storage"])
	}
}
