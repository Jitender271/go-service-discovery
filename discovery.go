package main

import (
	"log"

	"github.com/go-zookeeper/zk"
)

// discoverServices retrieves the list of services from Zookeeper and sets up watchers
func discoverServices(conn *zk.Conn, serviceType string) ([]string, error) {
	servicePath := "/services/" + serviceType

	// Check if the service type exists
	exists, _, err := conn.Exists(servicePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return []string{}, nil // No services registered under this type
	}

	// Get the initial list of services
	services, _, err := conn.Children(servicePath)
	if err != nil {
		return nil, err
	}

	log.Printf("Available services for %s: %v", serviceType, services)

	// Set a watcher to listen for changes
	go watchServices(conn, servicePath)

	return services, nil
}

// watchServices sets a watcher on the service node to listen for changes
func watchServices(conn *zk.Conn, servicePath string) {
	for {
		_, _, events, err := conn.ChildrenW(servicePath)
		if err != nil {
			log.Printf("Error watching services: %v", err)
			return
		}

		// Wait for an event
		event := <-events
		if event.Type == zk.EventNodeChildrenChanged {
			// Fetch the updated list of services
			services, _, err := conn.Children(servicePath)
			if err != nil {
				log.Printf("Error fetching updated services: %v", err)
				continue
			}
			log.Printf("Updated services for %s: %v", servicePath, services)
		}
	}
}
