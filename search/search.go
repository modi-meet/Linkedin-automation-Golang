package search

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// Run navigates, scrapes, and handles pagination
func Run(page *rod.Page, keyword string, limit int) []string {
	fmt.Println("Starting search for:", keyword)
	searchURL := fmt.Sprintf(
		"https://www.linkedin.com/search/results/people/?keywords=%s",
		url.QueryEscape(keyword),
	)

	// Use Navigate instead of MustNavigate to handle errors gracefully
	if err := page.Navigate(searchURL); err != nil {
		fmt.Printf("Error navigating to search URL: %v\n", err)
		return nil
	}

	// Wait for page to load, handle potential timeout/error
	if err := page.WaitStable(time.Second * 5); err != nil {
		fmt.Printf("Warning: Page did not stabilize: %v\n", err)
	}

	fmt.Println("Current URL:", page.MustInfo().URL)
	fmt.Println("Page Title:", page.MustInfo().Title)

	var allProfiles []string
	visited := make(map[string]bool)

	for len(allProfiles) < limit {
		// Scroll down to ensure results are loaded
		page.Mouse.Scroll(0, 1000, 5)
		time.Sleep(1 * time.Second)

		newProfiles := scrapeCurrentPage(page)

		// If no profiles found on current page, try scrolling more or waiting
		if len(newProfiles) == 0 {
			fmt.Println("No profiles found on current view, scrolling further...")
			page.Mouse.Scroll(0, 2000, 5)
			time.Sleep(2 * time.Second)
			newProfiles = scrapeCurrentPage(page)
		}

		for _, url := range newProfiles {
			if !visited[url] {
				visited[url] = true
				allProfiles = append(allProfiles, url)
				fmt.Println("   -> Collected:", url)
			}
			if len(allProfiles) >= limit {
				break
			}
		}

		if len(allProfiles) >= limit {
			break
		}

		fmt.Println("Looking for Next button...")

		page.Mouse.Scroll(0, 2000, 5)
		time.Sleep(1 * time.Second)

		nextBtn, err := page.Element("button[aria-label='Next']")

		if err != nil {
			fmt.Println("No 'Next' button found (or end of results). Stopping.")
			break
		}

		if disabled, _ := nextBtn.Attribute("disabled"); disabled != nil {
			fmt.Println("'Next' button is disabled. End of results.")
			break
		}

		nextBtn.MustClick()
		fmt.Println("⏳ Loading next page...")
		page.MustWaitStable()

		// sleep for human-like behavior
		time.Sleep(time.Duration(3+time.Now().Unix()%3) * time.Second)
	}

	return allProfiles
}

func scrapeCurrentPage(page *rod.Page) []string {
	var urls []string
	links := page.MustElements("a")

	for _, link := range links {
		hrefPtr, err := link.Attribute("href")
		if err != nil || hrefPtr == nil {
			continue
		}

		href := strings.Split(*hrefPtr, "?")[0]

		if strings.HasPrefix(href, "/in/") {
			href = "https://www.linkedin.com" + href
		}

		if strings.Contains(href, "linkedin.com/in/") {
			urls = append(urls, href)
		}
	}

	return urls
}
