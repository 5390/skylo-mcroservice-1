package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB(connStr string) *DB {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	return &DB{conn: conn}
}

func (db *DB) SaveFailedMessage(message string) error {
	_, err := db.conn.Exec("INSERT INTO failed_messages (message) VALUES ($1)", message)
	return err
}

func (db *DB) GetPendingMessages() ([]string, error) {
	rows, err := db.conn.Query("SELECT message FROM failed_messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var msg string
		if err := rows.Scan(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
