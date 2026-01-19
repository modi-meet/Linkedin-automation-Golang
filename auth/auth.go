package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/meetm/linkedin-automation-go/pkg/logger"
	"github.com/meetm/linkedin-automation-go/utils"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

var (
	ErrLoginFailed     = errors.New("login failed: could not verify successful login")
	ErrCredentialError = errors.New("login failed: invalid credentials or account issue")
	ErrCaptchaDetected = errors.New("login blocked: CAPTCHA or verification required")
)

func Login(page *rod.Page, log *logger.Logger) error {
	log.Printf("Navigating to login page...")

	var navErr error
	for i := 0; i < 3; i++ {
		navErr = page.Navigate("https://www.linkedin.com/login")
		if navErr == nil {
			break
		}
		log.Printf("Navigation attempt %d failed, retrying...", i+1)
		time.Sleep(2 * time.Second)
	}
	if navErr != nil {
		return fmt.Errorf("navigation failed: %w", navErr)
	}

	utils.LongRandomSleep(2, 4)
	page.MustWaitStable()

	if has, _, _ := page.Has("#username"); !has {
		log.Printf("Looking for alternate sign-in option...")
		if el, err := page.Timeout(5*time.Second).ElementR("button, a", "Sign in using another account"); err == nil {
			utils.HumanClick(page, el)
			utils.LongRandomSleep(1, 2)
		} else if el, err := page.Timeout(3*time.Second).ElementR("button, a", "Sign in"); err == nil {
			utils.HumanClick(page, el)
			utils.LongRandomSleep(1, 2)
		}
	}

	emailInput, err := utils.WaitForElement(page, "#username", 10*time.Second)
	if err != nil {
		return errors.New("could not find email input field")
	}

	email := os.Getenv("LINKEDIN_EMAIL")
	pass := os.Getenv("LINKEDIN_PASSWORD")

	if email == "" || pass == "" {
		return errors.New("LINKEDIN_EMAIL or LINKEDIN_PASSWORD not set in environment")
	}

	log.Printf("Entering credentials...")

	if err := utils.HumanType(page, emailInput, email); err != nil {
		return err
	}

	utils.RandomSleep(300, 600)

	passwordInput, err := utils.WaitForElement(page, "#password", 5*time.Second)
	if err != nil {
		return errors.New("could not find password input field")
	}

	if err := utils.HumanType(page, passwordInput, pass); err != nil {
		return err
	}

	utils.RandomSleep(500, 1000)

	log.Printf("Submitting login...")
	page.Keyboard.Press(input.Enter)

	utils.LongRandomSleep(3, 5)
	page.MustWaitStable()

	return validateLogin(page, log)
}

func validateLogin(page *rod.Page, log *logger.Logger) error {
	for attempt := 0; attempt < 40; attempt++ {
		currentURL := page.MustInfo().URL

		if strings.Contains(currentURL, "/checkpoint") || strings.Contains(currentURL, "/challenge") {
			if attempt == 0 {
				log.Printf("Security checkpoint detected - please solve manually...")
			}
			time.Sleep(3 * time.Second)
			continue
		}

		if strings.Contains(currentURL, "/login") {
			if _, err := page.Timeout(2 * time.Second).Element(".form__label--error"); err == nil {
				return ErrCredentialError
			}
			time.Sleep(2 * time.Second)
			continue
		}

		if strings.Contains(currentURL, "/feed") || strings.Contains(currentURL, "/mynetwork") {
			log.Printf("Login successful")
			return nil
		}

		log.Printf("Login completed, current URL: %s", currentURL)
		return nil
	}

	return ErrCaptchaDetected
}

func SaveCookies(browser *rod.Browser, filename string, log *logger.Logger) error {
	log.Printf("Saving cookies...")

	cookies, err := browser.GetCookies()
	if err != nil {
		return err
	}

	var params []*proto.NetworkCookieParam
	for _, c := range cookies {
		var expires proto.TimeSinceEpoch
		if c.Expires != 0 {
			expires = proto.TimeSinceEpoch(c.Expires)
		}

		param := &proto.NetworkCookieParam{
			Name:      c.Name,
			Value:     c.Value,
			Domain:    c.Domain,
			Path:      c.Path,
			Expires:   expires,
			Secure:    c.Secure,
			HTTPOnly:  c.HTTPOnly,
			SameSite:  c.SameSite,
			Priority:  c.Priority,
			SameParty: c.SameParty,
		}
		params = append(params, param)
	}

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadCookies(browser *rod.Browser, filename string, log *logger.Logger) error {
	log.Printf("Loading cookies from %s...", filename)

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("No cookie file found")
		return err
	}

	var params []*proto.NetworkCookieParam
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}

	if err := browser.SetCookies(params); err != nil {
		return err
	}

	log.Printf("Cookies loaded successfully")
	return nil
}
