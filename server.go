package main

import (
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

// Server holds the dependencies for the HTTP handlers
type Server struct {
	zkConn *zk.Conn
}

// NewServer initializes the Zookeeper connection and returns a Server instance
func NewServer() (*Server, error) {
	zkAddress := "localhost:2181" // You can make this configurable
	conn, _, err := zk.Connect([]string{zkAddress}, time.Second*5)
	if err != nil {
		return nil, err
	}

	// Ensure the base path exists
	basePath := "/services"
	exists, _, err := conn.Exists(basePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		_, err = conn.Create(basePath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return nil, err
		}
		log.Printf("Created base path: %s", basePath)
	}

	return &Server{zkConn: conn}, nil
}
