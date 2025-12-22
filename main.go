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

	u := launcher.New().
		Headless(false).
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)

	cookieFile := "linkedin_cookies.json"

	// load cookies
	err = auth.LoadCookies(browser, cookieFile)

	if err == nil {
		fmt.Println("Navigating to LinkedIn Feed...")
		page.MustNavigate("https://www.linkedin.com/feed/")

		if page.MustInfo().URL == "https://www.linkedin.com/login" {
			fmt.Println("Cookies expired. Logging in again.")
			auth.Login(page)
			auth.SaveCookies(browser, cookieFile)
		}
	} else {
		// login logic
		auth.Login(page)
		auth.SaveCookies(browser, cookieFile)
	}

	keyword := "Golang Developer"

	// search run for profiles
	profiles := search.Run(page, keyword, 3) // only 3 for testing

	// connection requests
	for _, profile := range profiles {
		msg := "Hi, I am a Go developer expanding my network. Would love to connect!"

		actions.SendConnectionRequest(page, profile, msg)

		fmt.Println("Cooling down b/w requests...")
		utils.RandomSleep(5000, 10000)
	}

	fmt.Println("Workflow Complete!!!")
}
