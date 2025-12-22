package main

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

func main() {
	fmt.Println("Booting up LinkedIn Bot...")

	u := launcher.New().
		Headless(false).
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)

	fmt.Println(" Verifying stealth status...")
	page.MustNavigate("https://bot.sannysoft.com")

	fmt.Println("Browser launched. Check the window for all 'Green' tests.")
	time.Sleep(30 * time.Second)
}
