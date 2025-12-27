package main

import (
	"fmt"
	"os"

	"github.com/meetm/linkedin-automation-go/actions"
	"github.com/meetm/linkedin-automation-go/auth"
	"github.com/meetm/linkedin-automation-go/search"
	"github.com/meetm/linkedin-automation-go/utils"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	"github.com/joho/godotenv"
)

func main() {
	// load env credentials
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	fmt.Println("Booting up LinkedIn Bot...")

	// Configure the browser launcher with stealth settings
	l := launcher.New().
		Headless(false)

	u := l.MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	// Create a stealth page
	page := stealth.MustPage(browser)

	cookieFile := "linkedin_cookies.json"

	// load cookies
	err = auth.LoadCookies(browser, cookieFile)

	if err == nil {
		fmt.Println("Session restored from cookies.")
	} else {
		// login logic
		fmt.Println("No valid session. Performing login...")
		auth.Login(page)
		if err := auth.SaveCookies(browser, cookieFile); err != nil {
			fmt.Println("Warning: Failed to save cookies:", err)
		}
	}

	keyword := "Software Developer"

	// search run for profiles
	profiles := search.Run(page, keyword, 3) // only 3 for testing

	if len(profiles) == 0 {
		fmt.Println("No profiles found. Exiting.")
		return
	}

	fmt.Printf("Found %d profiles. Starting connection requests...\n", len(profiles))

	// connection requests
	for _, profile := range profiles {
		msg := "Hi, I am a Go developer expanding my network. Would love to connect!"

		actions.SendConnectionRequest(page, profile, msg)

		fmt.Println("Cooling down b/w requests...")
		utils.RandomSleep(5000, 10000)
	}

	fmt.Println("Workflow Complete!!!")
}
