package main

import (
	"github.com/joho/godotenv"
	"github.com/meetm/linkedin-automation-go/api"
	"github.com/meetm/linkedin-automation-go/pkg/logger"
)

func main() {
	godotenv.Load()

	log := logger.New()

	server := api.NewServer(log)
	server.Start()
}
