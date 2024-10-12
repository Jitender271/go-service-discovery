package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the server with Zookeeper connection
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer server.zkConn.Close()

	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Register routes with appropriate handlers and methods
	router.HandleFunc("/register", server.handleRegister).Methods("POST")
	router.HandleFunc("/discover/{serviceType}", server.handleDiscover).Methods("GET")

	// Start the server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
