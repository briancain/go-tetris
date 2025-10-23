package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/briancain/go-tetris/internal/server/services"
	"github.com/briancain/go-tetris/internal/server/storage/memory"
)

func TestAuthHandler_UsernameConflicts(t *testing.T) {
	// Setup
	playerStore := memory.NewPlayerStore()
	authService := services.NewAuthService(playerStore)
	handler := NewAuthHandler(authService)

	t.Run("should return 409 for duplicate username", func(t *testing.T) {
		// First login request
		loginReq1 := LoginRequest{Username: "testuser"}
		body1, _ := json.Marshal(loginReq1)
		req1 := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body1))
		req1 = req1.WithContext(context.WithValue(req1.Context(), "requestID", "test-req-1"))
		w1 := httptest.NewRecorder()

		handler.Login(w1, req1)

		if w1.Code != http.StatusOK {
			t.Fatalf("Expected first login to succeed with status 200, got %d", w1.Code)
		}

		// Second login request with same username
		loginReq2 := LoginRequest{Username: "testuser"}
		body2, _ := json.Marshal(loginReq2)
		req2 := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body2))
		req2 = req2.WithContext(context.WithValue(req2.Context(), "requestID", "test-req-2"))
		w2 := httptest.NewRecorder()

		handler.Login(w2, req2)

		if w2.Code != http.StatusConflict {
			t.Errorf("Expected second login to fail with status 409, got %d", w2.Code)
		}

		expectedMessage := "Username is already in use. Please choose a different username."
		if w2.Body.String() != expectedMessage+"\n" {
			t.Errorf("Expected error message '%s', got '%s'", expectedMessage, w2.Body.String())
		}
	})

	t.Run("should allow login after logout", func(t *testing.T) {
		// First login
		loginReq := LoginRequest{Username: "reusable"}
		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), "requestID", "test-req-3"))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected login to succeed with status 200, got %d", w.Code)
		}

		// Parse response to get player ID
		var loginResp LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &loginResp)
		if err != nil {
			t.Fatalf("Failed to parse login response: %v", err)
		}

		// Logout
		logoutReq := httptest.NewRequest("POST", "/api/auth/logout", nil)
		logoutReq = logoutReq.WithContext(context.WithValue(logoutReq.Context(), "requestID", "test-req-4"))
		logoutReq = logoutReq.WithContext(context.WithValue(logoutReq.Context(), "playerID", loginResp.PlayerID))
		logoutW := httptest.NewRecorder()

		handler.Logout(logoutW, logoutReq)

		if logoutW.Code != http.StatusOK {
			t.Fatalf("Expected logout to succeed with status 200, got %d", logoutW.Code)
		}

		// Login again with same username should work
		req2 := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req2 = req2.WithContext(context.WithValue(req2.Context(), "requestID", "test-req-5"))
		w2 := httptest.NewRecorder()

		handler.Login(w2, req2)

		if w2.Code != http.StatusOK {
			t.Errorf("Expected login after logout to succeed with status 200, got %d", w2.Code)
		}
	})

	t.Run("should allow different usernames", func(t *testing.T) {
		// First user
		loginReq1 := LoginRequest{Username: "user1"}
		body1, _ := json.Marshal(loginReq1)
		req1 := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body1))
		req1 = req1.WithContext(context.WithValue(req1.Context(), "requestID", "test-req-6"))
		w1 := httptest.NewRecorder()

		handler.Login(w1, req1)

		if w1.Code != http.StatusOK {
			t.Fatalf("Expected first login to succeed with status 200, got %d", w1.Code)
		}

		// Second user with different username
		loginReq2 := LoginRequest{Username: "user2"}
		body2, _ := json.Marshal(loginReq2)
		req2 := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body2))
		req2 = req2.WithContext(context.WithValue(req2.Context(), "requestID", "test-req-7"))
		w2 := httptest.NewRecorder()

		handler.Login(w2, req2)

		if w2.Code != http.StatusOK {
			t.Errorf("Expected second login with different username to succeed with status 200, got %d", w2.Code)
		}
	})
}
