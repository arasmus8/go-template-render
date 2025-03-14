package server

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	port    int
	logFile *os.File
}

func NewServer(portVal int, logFile *os.File) *http.Server {
	NewServer := &Server{
		port:    portVal,
		logFile: logFile,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	return server
}
