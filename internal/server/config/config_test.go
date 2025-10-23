package config

import (
	"os"
	"testing"
)

func TestGetValue(t *testing.T) {
	tests := []struct {
		name         string
		flagValue    string
		envValue     string
		defaultValue string
		expected     string
	}{
		{"flag takes priority", "flag-val", "env-val", "default-val", "flag-val"},
		{"env fallback when no flag", "", "env-val", "default-val", "env-val"},
		{"default when neither", "", "", "default-val", "default-val"},
		{"empty flag uses env", "", "env-val", "default-val", "env-val"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env var if provided
			if tt.envValue != "" {
				os.Setenv("TEST_KEY", tt.envValue)
				defer os.Unsetenv("TEST_KEY")
			}

			result := getValue(tt.flagValue, "TEST_KEY", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoadWithDefaults(t *testing.T) {
	cfg, err := LoadWithFlags(false)
	if err != nil {
		t.Fatalf("LoadWithFlags() failed: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Port)
	}
	if cfg.RedisURL != "redis://localhost:6379" {
		t.Errorf("Expected default Redis URL, got %s", cfg.RedisURL)
	}
	if cfg.ServerURL != "http://localhost:8080" {
		t.Errorf("Expected default server URL, got %s", cfg.ServerURL)
	}
}

func TestLoadWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9000")
	os.Setenv("REDIS_URL", "redis://prod:6379")
	os.Setenv("SERVER_URL", "https://api.example.com")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("SERVER_URL")
	}()

	cfg, err := LoadWithFlags(false)
	if err != nil {
		t.Fatalf("LoadWithFlags() failed: %v", err)
	}

	if cfg.Port != "9000" {
		t.Errorf("Expected port 9000, got %s", cfg.Port)
	}
	if cfg.RedisURL != "redis://prod:6379" {
		t.Errorf("Expected prod Redis URL, got %s", cfg.RedisURL)
	}
	if cfg.ServerURL != "https://api.example.com" {
		t.Errorf("Expected prod server URL, got %s", cfg.ServerURL)
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		wantErr bool
	}{
		{"valid port", "8080", false},
		{"invalid port", "abc", true},
		{"empty port", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:      tt.port,
				RedisURL:  "redis://localhost:6379",
				ServerURL: "http://localhost:8080",
			}

			err := cfg.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
