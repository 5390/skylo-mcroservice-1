package config

import (
	"log"
	"os"
)

// Config holds application configurations.
type Config struct {
	DatabaseURL string
	ServerPort  string
}

// LoadConfig loads configurations from environment variables with defaults.
func LoadConfig() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/retry_db?sslmode=disable"),
		ServerPort:  getEnv("SERVER_PORT", "8081"),
	}
}

// getEnv fetches an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Environment variable %s not set, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}
