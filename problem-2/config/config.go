package config

import (
	"log"
	"os"
	"problem-2/db"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config holds all configuration values for the service.
type Config struct {
	InputFolder  string        // Path to the folder for input files
	OutputFolder string        // Path to the folder for output files
	DBHandler    *db.DBHandler // Database handler
	ScanInterval time.Duration // Time interval for folder scanning
}

// LoadConfig reads configurations and initializes a DB handler.
func LoadConfig() Config {
	// Retrieve the database URL from environment variables or use a default.
	dbURL := getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/retry_db?sslmode=disable")

	// Initialize the DB handler
	dbHandler, err := db.NewDBHandler(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize database handler: %v", err)
	}
	sqlFile, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = dbHandler.Conn.Exec(string(sqlFile))
	if err != nil {
		log.Fatal(err)
	}
	// Return the configuration with all necessary settings
	return Config{
		InputFolder:  getEnv("INPUT_FOLDER", "./inputs"),
		OutputFolder: getEnv("OUTPUT_FOLDER", "./outputs"),
		DBHandler:    dbHandler,
		ScanInterval: getEnvAsDuration("SCAN_INTERVAL", 30*time.Second),
	}
}

// getEnv fetches an environment variable or returns the default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Printf("Loaded %s from environment: %s", key, value)
		return value
	}
	log.Printf("Environment variable %s not set, using default: %s", key, defaultValue)
	return defaultValue
}

// getEnvAsDuration fetches an environment variable as a time.Duration or returns the default value.
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Invalid duration for %s: %s, using default: %v", key, valueStr, defaultValue)
		return defaultValue
	}

	return duration
}
