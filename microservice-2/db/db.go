package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// DB is a wrapper for the SQL connection.
type DB struct {
	Conn *sql.DB
}

// NewDB initializes a new database connection.
func NewDB(databaseURL string) *DB {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}

	log.Println("Connected to the database successfully.")
	return &DB{Conn: conn}
}

// InsertMessage inserts a received message into the database.
func (db *DB) InsertMessage(data string) error {
	_, err := db.Conn.Exec("INSERT INTO received_messages (data, received_at) VALUES ($1, CURRENT_TIMESTAMP)", data)
	return err
}
