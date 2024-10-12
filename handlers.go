package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Service represents a service to be registered
type Service struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Port    int    `json:"port" validate:"required"`
}

// handleRegister handles the /register POST endpoint
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var service Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate service fields (optional but recommended)
	if service.Name == "" || service.Address == "" || service.Port == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Register the service
	err := registerService(s.zkConn, service)
	if err != nil {
		log.Printf("Failed to register service: %v", err)
		http.Error(w, "Failed to register service", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service registered successfully"))
}

// handleDiscover handles the /discover/{serviceType} GET endpoint
func (s *Server) handleDiscover(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceType := vars["serviceType"]

	services, err := discoverServices(s.zkConn, serviceType)
	if err != nil {
		log.Printf("Failed to discover services: %v", err)
		http.Error(w, "Failed to discover services", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"services": services})
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
