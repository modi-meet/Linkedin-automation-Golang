package main

import (
	"fmt"

	"github.com/meetm/linkedin-automation-go/config"
)

func main() {
	cfg := config.Load()

	fmt.Println("Email:", cfg.LinkedInEmail)
}
