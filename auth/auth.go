package auth

import (
	"fmt"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
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
