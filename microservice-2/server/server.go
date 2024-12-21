package server

import (
	"encoding/json"
	"fmt"
	"log"
	"microservice-2/db"
	"net/http"
)

// Payload represents the structure of incoming data.
type Payload struct {
	Data string `json:"data"`
}

// Server holds dependencies for the HTTP server.
type Server struct {
	DB *db.DB
}

// NewServer initializes a new Server instance.
func NewServer(db *db.DB) *Server {
	return &Server{DB: db}
}

// Start runs the HTTP server on the specified port.
func (s *Server) Start(port string) {
	http.HandleFunc("/api/data", s.handleData)

	// CORS configuration
	http.HandleFunc("/", s.handleCORS)

	log.Printf("Starting Microservice-2 on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// handleCORS adds CORS headers to the response
func (s *Server) handleCORS(w http.ResponseWriter, r *http.Request) {
	// Allow all origins for testing (you can restrict to specific origins later)
	w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with your microservice-1 origin to restrict

	// Allow necessary methods
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	// Allow headers you expect to receive
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests for OPTIONS method
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Proceed to the next handler (the main logic handler)
}

// handleData handles POST requests to save received data.
func (s *Server) handleData(w http.ResponseWriter, r *http.Request) {
	// Log headers and body for debugging
	log.Printf("Request Headers: %v", r.Header)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Error decoding payload: %v", err) // Log the error
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received payload: %v", payload)

	if err := s.DB.InsertMessage(payload.Data); err != nil {
		fmt.Println("DB Error ::", err)
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	log.Printf("Saved data: %s", payload.Data)
	w.WriteHeader(http.StatusOK)
}
