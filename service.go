package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/go-zookeeper/zk"
)

// createPath ensures that the given path exists by creating any missing znodes
func createPath(conn *zk.Conn, path string) error {
	parts := strings.Split(path, "/")
	currentPath := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		currentPath += "/" + part
		exists, _, err := conn.Exists(currentPath)
		if err != nil {
			return err
		}
		if !exists {
			_, err := conn.Create(currentPath, []byte(""), 0, zk.WorldACL(zk.PermAll))
			if err != nil && err != zk.ErrNodeExists {
				return err
			}
			log.Printf("Created path: %s", currentPath)
		}
	}
	return nil
}

// registerService registers a new service in Zookeeper
func registerService(conn *zk.Conn, service Service) error {
	servicePath := "/services/" + service.Name

	// Ensure the service type path exists
	if err := createPath(conn, servicePath); err != nil {
		return err
	}

	// Create a unique node for the service instance using sequential ephemeral nodes
	nodePath := servicePath + "/service-"
	data := []byte(service.Address + ":" + strconv.Itoa(service.Port))
	createdPath, err := conn.CreateProtectedEphemeralSequential(nodePath, data, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}

	log.Printf("Service registered: %s at %s", createdPath, data)
	return nil
}
