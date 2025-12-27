package auth

import (
	"encoding/json"
	"os"

	"github.com/meetm/linkedin-automation-go/pkg/logger"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

func Login(page *rod.Page, log *logger.Logger) {
	log.Printf("Check: Are we logged in?")

	page.MustNavigate("https://www.linkedin.com/login")
	page.MustWaitStable()

	// Check if we are on the "Welcome Back" page (no username field)
	if has, _, _ := page.Has("#username"); !has {
		log.Printf("Username field not found. Checking for 'Sign in using another account'...")

		// Try to find the button/link
		el, err := page.ElementR("button, a", "Sign in using another account")
		if err == nil {
			log.Printf("Found 'Sign in using another account' button. Clicking...")
			el.MustClick()
			page.MustWaitStable()
		} else {
			// Fallback: maybe it's just "Sign in"
			el, err = page.ElementR("button, a", "Sign in")
			if err == nil {
				log.Printf("Found 'Sign in' button. Clicking...")
				el.MustClick()
				page.MustWaitStable()
			}
		}
	}

	emailInput := page.MustElement("#username")

	email := os.Getenv("LINKEDIN_EMAIL")
	pass := os.Getenv("LINKEDIN_PASS")

	log.Printf("Typing credentials...")

	emailInput.MustInput(email)
	page.MustElement("#password").MustInput(pass)

	log.Printf("Hitting Enter...")
	page.KeyActions().Press(input.Enter).MustDo()

	log.Printf("Waiting for home feed...")
	page.MustWaitStable()

	log.Printf("Login successful!")
}

// export browser cookies to a file
func SaveCookies(browser *rod.Browser, filename string, log *logger.Logger) error {
	log.Printf("Saving cookies to %s", filename)

	cookies, err := browser.GetCookies()
	if err != nil {
		return err
	}

	var params []*proto.NetworkCookieParam
	for _, c := range cookies {
		domain := c.Domain
		path := c.Path

		// optional fields
		var expires proto.TimeSinceEpoch
		if c.Expires != 0 {
			expires = proto.TimeSinceEpoch(c.Expires)
		}
		secure := c.Secure
		httpOnly := c.HTTPOnly
		sameParty := c.SameParty
		sameSite := c.SameSite
		priority := c.Priority

		param := &proto.NetworkCookieParam{
			Name:      c.Name,
			Value:     c.Value,
			Domain:    domain,
			Path:      path,
			Expires:   expires,
			Secure:    secure,
			HTTPOnly:  httpOnly,
			SameSite:  sameSite,
			Priority:  priority,
			SameParty: sameParty,
		}
		params = append(params, param)
	}

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// import cookies from a file if it exists
func LoadCookies(browser *rod.Browser, filename string, log *logger.Logger) error {
	log.Printf("Checking for cookie file: %s", filename)

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("No cookie file found. Fresh login required.")
		return err
	}

	var params []*proto.NetworkCookieParam
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}

	if err := browser.SetCookies(params); err != nil {
		return err
	}

	log.Printf("Cookies loaded! Session restored.")
	return nil
}
