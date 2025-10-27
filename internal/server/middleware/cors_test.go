package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseCORSOrigins(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single origin",
			input:    "https://example.com",
			expected: []string{"https://example.com"},
		},
		{
			name:     "multiple origins",
			input:    "https://example.com,http://localhost:3000,https://app.example.com",
			expected: []string{"https://example.com", "http://localhost:3000", "https://app.example.com"},
		},
		{
			name:     "origins with spaces",
			input:    "https://example.com, http://localhost:3000 , https://app.example.com",
			expected: []string{"https://example.com", "http://localhost:3000", "https://app.example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCORSOrigins(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseCORSOrigins(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCORS(t *testing.T) {
	allowedOrigins := []string{"https://example.com", "http://localhost:3000"}
	corsMiddleware := CORS(allowedOrigins)

	// Mock handler
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	tests := []struct {
		name           string
		method         string
		origin         string
		expectedStatus int
		expectedCORS   bool
	}{
		{
			name:           "allowed origin with GET",
			method:         "GET",
			origin:         "https://example.com",
			expectedStatus: http.StatusOK,
			expectedCORS:   true,
		},
		{
			name:           "allowed origin with OPTIONS",
			method:         "OPTIONS",
			origin:         "https://example.com",
			expectedStatus: http.StatusOK,
			expectedCORS:   true,
		},
		{
			name:           "disallowed origin with OPTIONS",
			method:         "OPTIONS",
			origin:         "https://malicious.com",
			expectedStatus: http.StatusForbidden,
			expectedCORS:   false,
		},
		{
			name:           "no origin header",
			method:         "GET",
			origin:         "",
			expectedStatus: http.StatusOK,
			expectedCORS:   false, // No CORS headers set for requests without Origin
		},
		{
			name:           "localhost origin",
			method:         "GET",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
			expectedCORS:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			rr := httptest.NewRecorder()
			handler := corsMiddleware(mockHandler)
			handler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			corsHeader := rr.Header().Get("Access-Control-Allow-Origin")
			if tt.expectedCORS {
				if corsHeader != tt.origin {
					t.Errorf("Expected CORS header %q, got %q", tt.origin, corsHeader)
				}
			} else {
				if corsHeader != "" {
					t.Errorf("Expected no CORS header, got %q", corsHeader)
				}
			}
		})
	}
}
