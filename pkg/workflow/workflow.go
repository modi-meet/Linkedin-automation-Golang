package workflow

import (
	"os"
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

type WorkflowStats struct {
	ProfilesFound   int
	RequestsSent    int
	RequestsSkipped int
	RequestsFailed  int
}

func Run(cfg Config, log *logger.Logger) {
	log.Printf("Starting LinkedIn automation...")

	browser, page, err := initBrowser(cfg.Headless, log)
	if err != nil {
		log.Printf("Browser initialization failed: %v", err)
		return
	}
	defer browser.MustClose()

	if cfg.Email != "" {
		os.Setenv("LINKEDIN_EMAIL", cfg.Email)
	}
	if cfg.Password != "" {
		os.Setenv("LINKEDIN_PASSWORD", cfg.Password)
	}

	log.Printf("Performing login...")
	if err := auth.Login(page, log); err != nil {
		log.Printf("Login failed: %v", err)
		return
	}

	utils.LongRandomSleep(2, 4)

	profiles := search.Run(page, cfg.Keyword, cfg.Limit, log)
	if len(profiles) == 0 {
		log.Printf("No profiles found. Exiting.")
		return
	}

	log.Printf("Found %d profiles. Starting connection requests...", len(profiles))

	stats := processProfiles(page, profiles, cfg.ConnectMessage, log)

	log.Printf("Workflow complete! Sent: %d, Skipped: %d, Failed: %d",
		stats.RequestsSent, stats.RequestsSkipped, stats.RequestsFailed)
}

func initBrowser(headless bool, log *logger.Logger) (*rod.Browser, *rod.Page, error) {
	userDataDir := getUserDataDir()
	log.Printf("Using browser profile: %s", userDataDir)

	l := launcher.New().
		Headless(headless).
		UserDataDir(userDataDir).
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-dev-shm-usage").
		Set("no-first-run").
		Set("no-default-browser-check").
		Set("disable-infobars").
		Set("disable-extensions").
		Set("disable-popup-blocking").
		Set("ignore-certificate-errors").
		Set("disable-background-networking").
		Set("disable-sync").
		Set("disable-translate").
		Set("metrics-recording-only").
		Set("safebrowsing-disable-auto-update").
		Set("password-store", "basic")

	log.Printf("Launching browser...")
	u, err := l.Launch()
	if err != nil {
		return nil, nil, err
	}

	browser := rod.New().ControlURL(u).MustConnect()
	time.Sleep(time.Second * 2)

	page := stealth.MustPage(browser)

	page.MustSetViewport(1920, 1080, 1, false)

	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		AcceptLanguage: "en-US,en;q=0.9",
		Platform:       "Win32",
	})

	applyStealthScripts(page)

	log.Printf("Browser ready")
	return browser, page, nil
}

func applyStealthScripts(page *rod.Page) {
	page.MustEval(`() => {
		Object.defineProperty(navigator, 'webdriver', {get: () => undefined});
		
		Object.defineProperty(navigator, 'languages', {get: () => ['en-US', 'en']});
		
		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5]
		});
		
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({ state: Notification.permission }) :
				originalQuery(parameters)
		);
		
		Object.defineProperty(navigator, 'platform', {get: () => 'Win32'});
		
		Object.defineProperty(navigator, 'hardwareConcurrency', {get: () => 8});
		
		Object.defineProperty(navigator, 'deviceMemory', {get: () => 8});
		
		const getParameter = WebGLRenderingContext.prototype.getParameter;
		WebGLRenderingContext.prototype.getParameter = function(parameter) {
			if (parameter === 37445) return 'Intel Inc.';
			if (parameter === 37446) return 'Intel Iris OpenGL Engine';
			return getParameter.apply(this, arguments);
		};
		
		window.chrome = {runtime: {}};
	}`)
}

func getUserDataDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.linkedin-automation-profile"
}

func processProfiles(page *rod.Page, profiles []string, message string, log *logger.Logger) WorkflowStats {
	stats := WorkflowStats{ProfilesFound: len(profiles)}

	if message == "" {
		message = "Hi, I am a Go developer expanding my network. Would love to connect!"
	}

	for i, profile := range profiles {
		log.Printf("Processing %d/%d...", i+1, len(profiles))

		result := actions.SendConnectionRequest(page, profile, message, log)

		if result.Success {
			stats.RequestsSent++
		} else if result.Skipped {
			stats.RequestsSkipped++
		} else {
			stats.RequestsFailed++
		}

		if i < len(profiles)-1 {
			log.Printf("Cooling down...")
			utils.LongRandomSleep(5, 12)
		}
	}

	return stats
}
