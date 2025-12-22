package main

import (
	"fmt"
	"os"
	"time"

	"github.com/meetm/linkedin-automation-go/auth"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	fmt.Println("Booting up LinkedIn Bot...")

	u := launcher.New().
		Headless(false).
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)

	// Login Logic
	auth.Login(page)

	fmt.Println("👀Check the browser: Are you on the Feed?")
	time.Sleep(1 * time.Minute)
}
