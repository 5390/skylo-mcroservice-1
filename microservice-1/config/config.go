package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config represents the overall configuration for Microservice-1.
type Config struct {
	QueueConfig QueueConfig
	RetryConfig RetryConfig
	DatabaseURL string
}

// QueueConfig holds configurations for the message queue.
type QueueConfig struct {
	Broker  string
	Topic   string
	GroupID string
}

// RetryConfig holds configurations for the retry mechanism.
type RetryConfig struct {
	TargetURL  string
	RetryDelay time.Duration
}

// LoadConfig reads the configuration from environment variables and provides defaults.
func LoadConfig() Config {
	return Config{
		QueueConfig: QueueConfig{
			Broker:  getEnv("QUEUE_BROKER", "localhost:9092"),
			Topic:   getEnv("QUEUE_TOPIC", "my-topic"),
			GroupID: getEnv("QUEUE_GROUP_ID", "my-group"),
		},
		RetryConfig: RetryConfig{
			TargetURL:  getEnv("RETRY_TARGET_URL", "http://microservice-2:8081/api/data"),
			RetryDelay: getEnvAsDuration("RETRY_DELAY", 10*time.Second),
		},
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/retry_db?sslmode=disable"),
	}
}

// getEnv reads an environment variable or returns the default value.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Environment variable %s not set, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvAsDuration reads an environment variable as a duration or returns the default value.
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Invalid duration for %s: %s, using default: %v", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvAsInt reads an environment variable as an integer or returns the default value.
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid integer for %s: %s, using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}
