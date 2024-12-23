package db

import (
	"database/sql"
	"fmt"
	"log"
	"problem-2/models"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DBHandler wraps the SQL DB connection and provides methods for database operations.
type DBHandler struct {
	Conn *sql.DB
}

// NewDBHandler initializes a new DBHandler with a connection to the database.
func NewDBHandler(dbURL string) (*DBHandler, error) {
	// Establish connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	log.Println("Connected to the database successfully.")
	return &DBHandler{Conn: db}, nil
}

// InsertSIMRecord inserts a SIM record into the database.
func (db *DBHandler) InsertSIMRecord(record models.SIMRecord) error {
	_, err := db.Conn.Exec(
		`INSERT INTO sim_records (imsi, pin1, puk1, pin2, puk2, aam1, ki_umts_enc, acc)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		record.IMSI, record.PIN1, record.PUK1, record.PIN2, record.PUK2, record.AAM1, record.KIUMTSEnc, record.ACC,
	)
	if err != nil {
		return fmt.Errorf("failed to insert SIM record: %w", err)
	}
	return nil
}

// IMSIExists checks if a given IMSI already exists in the database.
func (db *DBHandler) IMSIExists(imsi string) (bool, error) {
	var exists bool
	err := db.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM sim_records WHERE imsi = $1)", imsi).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check IMSI existence: %w", err)
	}
	return exists, nil
}
