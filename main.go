package main

import (
	"github.com/meetm/linkedin-automation-go/api"
	"github.com/meetm/linkedin-automation-go/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Initialize and start API server
	server := api.NewServer(log)
	server.Start()
}
