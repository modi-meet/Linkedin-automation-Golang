package actions

import (
	"errors"
	"time"

	"github.com/meetm/linkedin-automation-go/pkg/logger"
	"github.com/meetm/linkedin-automation-go/utils"

	"github.com/go-rod/rod"
)

var (
	ErrAlreadyConnected = errors.New("already connected or pending")
	ErrFollowOnly       = errors.New("profile only allows follow")
	ErrRateLimited      = errors.New("rate limited by LinkedIn")
	ErrConnectFailed    = errors.New("failed to send connection request")
)

type ConnectionResult struct {
	ProfileURL string
	Success    bool
	Error      error
	Skipped    bool
	Reason     string
}

func SendConnectionRequest(page *rod.Page, profileURL, message string, log *logger.Logger) ConnectionResult {
	result := ConnectionResult{ProfileURL: profileURL}

	log.Printf("Visiting: %s", profileURL)

	if err := page.Navigate(profileURL); err != nil {
		result.Error = err
		log.Printf("Navigation failed: %v", err)
		return result
	}

	utils.LongRandomSleep(2, 4)
	page.MustWaitStable()

	if isAlreadyConnected(page) {
		result.Skipped = true
		result.Reason = "already connected"
		log.Printf("Skipping: already connected")
		return result
	}

	if isPending(page) {
		result.Skipped = true
		result.Reason = "pending request"
		log.Printf("Skipping: pending request exists")
		return result
	}

	connectBtn := findConnectButton(page)
	if connectBtn == nil {
		connectBtn = findConnectInMore(page, log)
	}

	if connectBtn == nil {
		result.Skipped = true
		result.Reason = "no connect option"
		result.Error = ErrFollowOnly
		log.Printf("Skipping: no connect button found")
		return result
	}

	log.Printf("Clicking connect...")
	if err := utils.HumanClick(page, connectBtn); err != nil {
		result.Error = err
		return result
	}

	utils.RandomSleep(800, 1500)
	page.MustWaitStable()

	if !handleConnectionModal(page, message, log) {
		if sendWithoutNote(page, log) {
			result.Success = true
			log.Printf("Request sent (without note)")
			return result
		}
		result.Error = ErrConnectFailed
		log.Printf("Failed to complete connection flow")
		return result
	}

	result.Success = true
	log.Printf("Request sent successfully")
	return result
}

func isAlreadyConnected(page *rod.Page) bool {
	el, err := page.Timeout(2*time.Second).ElementR("button", "Message")
	if err != nil {
		return false
	}

	if el2, _ := page.Timeout(1*time.Second).ElementR("button", "Connect"); el2 != nil {
		return false
	}

	return utils.IsElementVisible(el)
}

func isPending(page *rod.Page) bool {
	el, err := page.Timeout(2*time.Second).ElementR("button", "Pending")
	if err != nil {
		return false
	}
	return utils.IsElementVisible(el)
}

func findConnectButton(page *rod.Page) *rod.Element {
	el, err := page.Timeout(3*time.Second).ElementR("button", "Connect")
	if err != nil {
		return nil
	}
	if !utils.IsElementVisible(el) {
		return nil
	}
	return el
}

func findConnectInMore(page *rod.Page, log *logger.Logger) *rod.Element {
	moreBtn, err := page.Timeout(3 * time.Second).Element("[aria-label='More actions']")
	if err != nil {
		return nil
	}

	log.Printf("Checking more actions menu...")
	if err := utils.HumanClick(page, moreBtn); err != nil {
		return nil
	}

	utils.RandomSleep(500, 1000)

	connectOption, err := page.Timeout(3*time.Second).ElementR("div[role='button'], span", "Connect")
	if err != nil {
		page.Keyboard.Press('\x1b')
		return nil
	}

	return connectOption
}

func handleConnectionModal(page *rod.Page, message string, log *logger.Logger) bool {
	addNoteBtn, err := page.Timeout(3*time.Second).ElementR("button", "Add a note")
	if err != nil {
		return false
	}

	if err := utils.HumanClick(page, addNoteBtn); err != nil {
		return false
	}

	utils.RandomSleep(500, 1000)

	textarea, err := utils.WaitForElement(page, "textarea[name='message']", 5*time.Second)
	if err != nil {
		textarea, err = utils.WaitForElement(page, "textarea", 3*time.Second)
		if err != nil {
			return false
		}
	}

	log.Printf("Typing message...")
	if err := utils.HumanType(page, textarea, message); err != nil {
		return false
	}

	utils.RandomSleep(500, 1000)

	sendBtn, err := page.Timeout(3*time.Second).ElementR("button", "Send")
	if err != nil {
		return false
	}

	if err := utils.HumanClick(page, sendBtn); err != nil {
		return false
	}

	utils.RandomSleep(500, 1000)
	return true
}

func sendWithoutNote(page *rod.Page, log *logger.Logger) bool {
	sendBtn, err := page.Timeout(2*time.Second).ElementR("button", "Send without a note")
	if err != nil {
		sendBtn, err = page.Timeout(2*time.Second).ElementR("button", "Send")
		if err != nil {
			return false
		}
	}

	log.Printf("Sending without note as fallback...")
	if err := utils.HumanClick(page, sendBtn); err != nil {
		return false
	}

	utils.RandomSleep(500, 1000)
	return true
}
