package main

import (
	"fmt"
	"os"
	"time"

	"github.com/meetm/linkedin-automation-go/auth"
	"github.com/meetm/linkedin-automation-go/search"

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
	profiles := search.Run(page, keyword, 50)

	fmt.Println("Collected total:", len(profiles))
	time.Sleep(1 * time.Minute)
}
