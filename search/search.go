package search

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/meetm/linkedin-automation-go/pkg/logger"
	"github.com/meetm/linkedin-automation-go/utils"

	"github.com/go-rod/rod"
)

func Run(page *rod.Page, keyword string, limit int, log *logger.Logger) []string {
	log.Printf("Searching for: %s", keyword)

	searchURL := fmt.Sprintf(
		"https://www.linkedin.com/search/results/people/?keywords=%s",
		url.QueryEscape(keyword),
	)

	if err := page.Navigate(searchURL); err != nil {
		log.Printf("Navigation error: %v", err)
		return nil
	}

	utils.LongRandomSleep(3, 5)

	if err := page.WaitStable(time.Second * 5); err != nil {
		log.Printf("Page stability warning: %v", err)
	}

	log.Printf("Search page loaded: %s", page.MustInfo().URL)

	if hasNoResults(page) {
		log.Printf("No search results found")
		return nil
	}

	var allProfiles []string
	visited := make(map[string]bool)
	pageNum := 1

	for len(allProfiles) < limit {
		log.Printf("Scanning page %d...", pageNum)

		utils.HumanScroll(page, 300)
		utils.RandomSleep(500, 1000)

		profiles := scrapeCurrentPage(page)

		if len(profiles) == 0 {
			log.Printf("No profiles found, scrolling further...")
			utils.HumanScroll(page, 800)
			utils.LongRandomSleep(1, 2)
			profiles = scrapeCurrentPage(page)
		}

		for _, profileURL := range profiles {
			if visited[profileURL] {
				continue
			}
			visited[profileURL] = true
			allProfiles = append(allProfiles, profileURL)
			log.Printf("Found: %s", profileURL)

			if len(allProfiles) >= limit {
				break
			}
		}

		if len(allProfiles) >= limit {
			break
		}

		utils.HumanScroll(page, 1500)
		utils.RandomSleep(800, 1500)

		if !goToNextPage(page, log) {
			break
		}

		pageNum++
		utils.LongRandomSleep(2, 4)
	}

	log.Printf("Collected %d profiles", len(allProfiles))
	return allProfiles
}

func hasNoResults(page *rod.Page) bool {
	el, err := page.Timeout(3*time.Second).ElementR("div", "No results found")
	if err != nil {
		return false
	}
	return utils.IsElementVisible(el)
}

func goToNextPage(page *rod.Page, log *logger.Logger) bool {
	nextBtn, err := page.Timeout(5 * time.Second).Element("button[aria-label='Next']")
	if err != nil {
		log.Printf("No next page available")
		return false
	}

	disabled, _ := nextBtn.Attribute("disabled")
	if disabled != nil {
		log.Printf("Reached last page")
		return false
	}

	if err := utils.HumanClick(page, nextBtn); err != nil {
		log.Printf("Failed to click next: %v", err)
		return false
	}

	log.Printf("Loading next page...")
	page.MustWaitStable()
	return true
}

func scrapeCurrentPage(page *rod.Page) []string {
	var urls []string
	links, err := page.Elements("a")
	if err != nil {
		return urls
	}

	for _, link := range links {
		hrefPtr, err := link.Attribute("href")
		if err != nil || hrefPtr == nil {
			continue
		}

		href := strings.Split(*hrefPtr, "?")[0]

		if strings.HasPrefix(href, "/in/") {
			href = "https://www.linkedin.com" + href
		}

		if strings.Contains(href, "linkedin.com/in/") && !strings.Contains(href, "/search/") {
			urls = append(urls, href)
		}
	}

	return urls
}
