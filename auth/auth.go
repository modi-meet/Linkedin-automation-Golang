package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

func Login(page *rod.Page) {
	fmt.Println("Check: Are we logged in?")

	page.MustNavigate("https://www.linkedin.com/login")

	emailInput := page.MustWaitStable().MustElement("#username")

	email := os.Getenv("LINKEDIN_EMAIL")
	pass := os.Getenv("LINKEDIN_PASS")

	fmt.Println("Typing credentials...")

	emailInput.MustInput(email)
	page.MustElement("#password").MustInput(pass)

	fmt.Println("Hitting Enter...")
	page.KeyActions().Press(input.Enter).MustDo()

	fmt.Println("Waiting for home feed...")
	page.MustWaitStable()

	fmt.Println("Login successful!")
}

// export browser cookies to a file
func SaveCookies(browser *rod.Browser, filename string) error {
	fmt.Println("Saving cookies to", filename)

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
func LoadCookies(browser *rod.Browser, filename string) error {
	fmt.Println("Checking for cookie file:", filename)

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("No cookie file found. Fresh login required.")
		return err
	}

	var params []*proto.NetworkCookieParam
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}

	if err := browser.SetCookies(params); err != nil {
		return err
	}

	fmt.Println("Cookies loaded! Session restored.")
	return nil
}
