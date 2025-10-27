package tetris

import (
	"runtime"
	"testing"
)

func TestGetServerURL(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "default localhost",
			expected: "http://localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getServerURL()

			// For non-WebAssembly builds, should always return localhost
			if runtime.GOARCH != "wasm" {
				if result != "http://localhost:8080" {
					t.Errorf("getServerURL() = %q, want %q", result, "http://localhost:8080")
				}
			} else {
				// For WebAssembly builds, we can't easily test JavaScript interaction
				// in unit tests, so we just verify it doesn't panic
				if result == "" {
					t.Error("getServerURL() returned empty string")
				}
			}
		})
	}
}

func TestGetServerURL_NonWasm(t *testing.T) {
	// This test specifically verifies the non-WebAssembly path
	if runtime.GOARCH == "wasm" {
		t.Skip("Skipping non-WebAssembly test on WebAssembly build")
	}

	result := getServerURL()
	expected := "http://localhost:8080"

	if result != expected {
		t.Errorf("getServerURL() = %q, want %q", result, expected)
	}
}
