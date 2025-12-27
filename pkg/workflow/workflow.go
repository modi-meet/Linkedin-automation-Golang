package workflow

import (
	"os"

	"github.com/meetm/linkedin-automation-go/actions"
	"github.com/meetm/linkedin-automation-go/auth"
	"github.com/meetm/linkedin-automation-go/pkg/logger"
	"github.com/meetm/linkedin-automation-go/search"
	"github.com/meetm/linkedin-automation-go/utils"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

type Config struct {
	Email          string
	Password       string
	Keyword        string
	Limit          int
	ConnectMessage string
	Headless       bool
}

func Run(cfg Config, log *logger.Logger) {
	log.Printf("Booting up LinkedIn Bot...")

	l := launcher.New().
		Headless(cfg.Headless)

	u := l.MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)

	cookieFile := "linkedin_cookies.json"

	if cfg.Email != "" {
		os.Setenv("LINKEDIN_EMAIL", cfg.Email)
	}
	if cfg.Password != "" {
		os.Setenv("LINKEDIN_PASSWORD", cfg.Password)
	}

	// load cookies
	err := auth.LoadCookies(browser, cookieFile, log)

	if err == nil {
		log.Printf("Session restored from cookies.")
	} else {
		// login logic
		log.Printf("No valid session. Performing login...")
		auth.Login(page, log)
		if err := auth.SaveCookies(browser, cookieFile, log); err != nil {
			log.Printf("Warning: Failed to save cookies: %v", err)
		}
	}

	// search run for profiles
	profiles := search.Run(page, cfg.Keyword, cfg.Limit, log)

	if len(profiles) == 0 {
		log.Printf("No profiles found. Exiting.")
		return
	}

	log.Printf("Found %d profiles. Starting connection requests...", len(profiles))

	// connection requests
	for _, profile := range profiles {
		msg := cfg.ConnectMessage
		if msg == "" {
			msg = "Hi, I am a Go developer expanding my network. Would love to connect!"
		}

		actions.SendConnectionRequest(page, profile, msg, log)

		log.Printf("Cooling down b/w requests...")
		utils.RandomSleep(5000, 10000)
	}

	log.Printf("Workflow Complete!!!")
}
