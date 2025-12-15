package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LinkedInEmail    string
	LinkedInPassword string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	email := os.Getenv("LINKEDIN_EMAIL")
	password := os.Getenv("LINKEDIN_PASSWORD")

	if email == "" || password == "" {
		log.Fatal("LINKEDIN_EMAIL or LINKEDIN_PASSWORD not set")
	}

	return &Config{
		LinkedInEmail:    email,
		LinkedInPassword: password,
	}
}
