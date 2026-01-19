package actions

import (
	"github.com/meetm/linkedin-automation-go/pkg/logger"
	"github.com/meetm/linkedin-automation-go/utils"

	"github.com/go-rod/rod"
)

func SendConnectionRequest(page *rod.Page, profileURL, message string, log *logger.Logger) {
	log.Printf("Visiting: %s", profileURL)

	// Navigate with error handling
	if err := page.Navigate(profileURL); err != nil {
		log.Printf("Error navigating to profile: %v", err)
		return
	}

	page.MustWaitStable()
	utils.RandomSleep(2000, 4000)

	connectBtn, err := page.ElementR("button", "Connect")

	if err != nil {
		log.Printf("'Connect' button not visible. Checking 'More'...")

		moreBtn, err := page.Element("[aria-label='More actions']")
		if err != nil {
			log.Printf("Could not find 'Connect' or 'More' button. Skipping.")
			return
		}

		utils.HumanClick(page, moreBtn)
		utils.RandomSleep(500, 1000)

		connectBtn, err = page.ElementR("div[role='button']", "Connect")
		if err != nil {
			log.Printf("'Connect' option not found in dropdown. (Might be Follow only).")
			return
		}
	}

	log.Printf("Found Connect button. Clicking...")
	utils.HumanClick(page, connectBtn)

	page.MustWaitStable()

	// handle "Add a note" button
	addNoteBtn, err := page.ElementR("button", "Add a note")
	if err == nil {
		utils.HumanClick(page, addNoteBtn)
		utils.RandomSleep(500, 1000)

		log.Printf("Typing message...")
		textArea := page.MustElement("textarea[name='message']")
		textArea.MustInput(message)
		utils.RandomSleep(500, 1000)

		sendBtn, err := page.ElementR("button", "Send")
		if err == nil {
			utils.HumanClick(page, sendBtn)
			log.Printf("📨 Request Sent!")
		}
	} else {
		log.Printf("'Add a note' button not found. Aborting to be safe.")
	}
}
