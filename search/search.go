package search

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-rod/rod"
)

func Run(page *rod.Page, keyword string) []string {
	fmt.Println("Starting search for:", keyword)

	searchURL := fmt.Sprintf(
		"https://www.linkedin.com/search/results/people/?keywords=%s",
		url.QueryEscape(keyword),
	)

	page.MustNavigate(searchURL)
	page.MustWaitStable()
	fmt.Println("Scraping profile URLs...")

	links := page.MustElements("a")

	seen := map[string]struct{}{}
	var profileURLs []string

	for _, link := range links {
		hrefPtr, err := link.Attribute("href")
		if err != nil || hrefPtr == nil {
			continue
		}

		href := strings.Split(*hrefPtr, "?")[0]

		// Normalize relative profile links to absolute
		if strings.HasPrefix(href, "/in/") {
			href = "https://www.linkedin.com" + href
		}

		// Keep only actual LinkedIn profile links
		if !strings.HasPrefix(href, "https://www.linkedin.com/in/") && !strings.HasPrefix(href, "https://linkedin.com/in/") {
			continue
		}

		if _, ok := seen[href]; ok {
			continue
		}
		seen[href] = struct{}{}
		profileURLs = append(profileURLs, href)
	}

	fmt.Printf("Found %d profiles on this page.\n", len(profileURLs))
	for _, u := range profileURLs {
		fmt.Println("   ->", u)
	}

	return profileURLs
}
