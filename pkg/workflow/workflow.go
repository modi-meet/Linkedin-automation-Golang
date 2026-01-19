package workflow

import (
	"os"

	"strings"
	"time"

	"github.com/meetm/linkedin-automation-go/actions"
	"github.com/meetm/linkedin-automation-go/auth"
	"github.com/meetm/linkedin-automation-go/pkg/logger"
	"github.com/meetm/linkedin-automation-go/search"
	"github.com/meetm/linkedin-automation-go/utils"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
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
		Headless(cfg.Headless).
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-infobars", "true").
		Set("start-maximized").
		Set("disable-web-security")

	log.Printf("Launching browser...")
	u := l.MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	log.Printf("Browser connected. Initializing page...")
	time.Sleep(time.Second * 2)

	page := stealth.MustPage(browser)
	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})

	log.Printf("Page initialized.")

	cookieFile := "linkedin_cookies.json"

	if cfg.Email != "" {
		os.Setenv("LINKEDIN_EMAIL", cfg.Email)
	}
	if cfg.Password != "" {
		os.Setenv("LINKEDIN_PASSWORD", cfg.Password)
	}

	// load cookies
	err := auth.LoadCookies(browser, cookieFile, log)
	sessionValid := false

	if err == nil {
		log.Printf("Cookies loaded! Session restored. Verifying...")
		// Navigate to feed to check if cookies are actually good
		// We use a retry loop here because sometimes the first request fails
		for i := 0; i < 3; i++ {
			if err := page.Navigate("https://www.linkedin.com/feed/"); err == nil {
				break
			}
			log.Printf("Navigation failed, retrying (%d/3)...", i+1)
			time.Sleep(time.Second * 2)
		}

		page.WaitStable(time.Second * 5)

		// Check if we are actually on the feed or redirected to login
		url := page.MustInfo().URL
		if strings.Contains(url, "linkedin.com/feed") {
			sessionValid = true
			log.Printf("Session verified.")
		} else {
			log.Printf("Session invalid (redirected to %s).", url)
		}
	}

	if !sessionValid {
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
