package main

import (
	"fmt"
	"log"
	"microservice-2/config"
	"microservice-2/db"
	"microservice-2/server"
	"os"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the database
	database := db.NewDB(cfg.DatabaseURL)
	defer database.Conn.Close()
	// Create the table if not already created
	sqlFile, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Conn.Exec(string(sqlFile))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table created successfully!")

	// Start the HTTP server
	srv := server.NewServer(database)
	srv.Start(cfg.ServerPort)
}
