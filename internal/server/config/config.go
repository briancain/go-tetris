package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        string
	RedisURL    string
	ServerURL   string
	CORSOrigins string
}

func Load() (*Config, error) {
	return LoadWithFlags(true)
}

func LoadWithFlags(parseFlags bool) (*Config, error) {
	var port, redisURL, serverURL, corsOrigins string

	if parseFlags && !flag.Parsed() {
		portFlag := flag.String("port", "", "Server port")
		redisURLFlag := flag.String("redis-url", "", "Redis connection URL")
		serverURLFlag := flag.String("server-url", "", "Public server URL")
		corsOriginsFlag := flag.String("cors-origins", "", "Comma-separated list of allowed CORS origins")
		flag.Parse()

		port = *portFlag
		redisURL = *redisURLFlag
		serverURL = *serverURLFlag
		corsOrigins = *corsOriginsFlag
	}

	cfg := &Config{
		Port:        getValue(port, "PORT", "8080"),
		RedisURL:    getValue(redisURL, "REDIS_URL", ""),
		ServerURL:   getValue(serverURL, "SERVER_URL", "http://localhost:8080"),
		CORSOrigins: getValue(corsOrigins, "CORS_ORIGINS", "http://localhost:3000,http://localhost:8080"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if c.ServerURL == "" {
		return fmt.Errorf("SERVER_URL is required")
	}

	// Validate port is numeric
	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("PORT must be numeric: %s", c.Port)
	}

	return nil
}

// getValue returns CLI flag value, then env var, then default
func getValue(flagValue, envKey, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	return defaultValue
}
