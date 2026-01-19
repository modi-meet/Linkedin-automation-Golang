package utils

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// init seeds randomness using crypto/rand for better unpredictability
func cryptoRandInt(min, max int) int {
	if max <= min {
		return min
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		return min
	}
	return int(n.Int64()) + min
}

func cryptoRandFloat() float64 {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return 0.5
	}
	return float64(n.Int64()) / 1000000.0
}

func HumanType(page *rod.Page, el *rod.Element, text string) error {
	if el == nil {
		return nil
	}

	if err := el.Focus(); err != nil {
		return err
	}

	RandomSleep(200, 400)

	for _, char := range text {
		// Type single character
		if err := el.Input(string(char)); err != nil {
			return err
		}

		// Random delay between keystrokes (50-200ms for normal chars)
		// Longer delays for spaces and punctuation (like real typing)
		if char == ' ' || char == '.' || char == ',' || char == '!' || char == '?' {
			RandomSleep(100, 300) // Slightly longer for punctuation/space
		} else {
			RandomSleep(50, 180) // Normal character delay
		}
	}

	RandomSleep(200, 500) // Small pause after typing
	return nil
}

// HumanClick moves mouse to element with human-like movement and clicks
func HumanClick(page *rod.Page, el *rod.Element) error {
	if el == nil {
		return nil
	}

	shape, err := el.Shape()
	if err != nil {
		return err
	}

	box := shape.Box()
	targetX := box.X + (box.Width / 2)
	targetY := box.Y + (box.Height / 2)

	// Add jitter to not always click exact center
	jitter := 8.0
	targetX += (cryptoRandFloat() * jitter) - (jitter / 2)
	targetY += (cryptoRandFloat() * jitter) - (jitter / 2)

	// Move mouse with natural curve
	moveMouseWithBezier(page, targetX, targetY)

	// Slight pause before clicking (humans don't click instantly)
	RandomSleep(80, 200)

	// Click
	if err := page.Mouse.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	// Small pause after clicking
	RandomSleep(150, 400)

	return nil
}

// moveMouseWithBezier moves mouse along a bezier curve for natural-looking movement
func moveMouseWithBezier(page *rod.Page, targetX, targetY float64) {
	// Get current mouse position (approximate from page center if unknown)
	startX := cryptoRandFloat() * 100
	startY := cryptoRandFloat() * 100

	// Generate random control points for bezier curve
	// This creates a natural arc rather than straight line
	ctrlX1 := startX + (targetX-startX)*0.3 + (cryptoRandFloat()*100 - 50)
	ctrlY1 := startY + (targetY-startY)*0.3 + (cryptoRandFloat()*100 - 50)
	ctrlX2 := startX + (targetX-startX)*0.7 + (cryptoRandFloat()*60 - 30)
	ctrlY2 := startY + (targetY-startY)*0.7 + (cryptoRandFloat()*60 - 30)

	// Number of steps for smooth movement
	steps := 25 + cryptoRandInt(0, 15)

	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)

		// Cubic bezier formula
		x := cubicBezier(startX, ctrlX1, ctrlX2, targetX, t)
		y := cubicBezier(startY, ctrlY1, ctrlY2, targetY, t)

		page.Mouse.MoveTo(proto.NewPoint(x, y))

		// Variable delay between movements (faster in middle, slower at ends)
		baseDelay := 5 + cryptoRandInt(0, 8)
		// Ease in/out effect
		if t < 0.2 || t > 0.8 {
			baseDelay += cryptoRandInt(2, 6)
		}
		time.Sleep(time.Duration(baseDelay) * time.Millisecond)
	}

	// Final move to exact target
	page.Mouse.MoveTo(proto.NewPoint(targetX, targetY))
}

// cubicBezier calculates point on cubic bezier curve
func cubicBezier(p0, p1, p2, p3, t float64) float64 {
	return math.Pow(1-t, 3)*p0 +
		3*math.Pow(1-t, 2)*t*p1 +
		3*(1-t)*math.Pow(t, 2)*p2 +
		math.Pow(t, 3)*p3
}

// HumanScroll scrolls the page gradually like a human would
func HumanScroll(page *rod.Page, totalDistance int) {
	if totalDistance == 0 {
		return
	}

	direction := 1
	if totalDistance < 0 {
		direction = -1
		totalDistance = -totalDistance
	}

	scrolled := 0
	for scrolled < totalDistance {
		// Random scroll chunk (50-200 pixels at a time)
		chunk := cryptoRandInt(50, 200)
		if scrolled+chunk > totalDistance {
			chunk = totalDistance - scrolled
		}

		page.Mouse.Scroll(0, float64(chunk*direction), 1)
		scrolled += chunk

		// Variable pause between scroll chunks (100-400ms)
		// Humans scroll in bursts, not continuously
		RandomSleep(100, 400)

		// Occasionally pause longer (simulating reading)
		if cryptoRandInt(0, 10) < 2 {
			RandomSleep(500, 1500)
		}
	}
}

// RandomSleep pauses for a random duration between min and max milliseconds
func RandomSleep(min, max int) {
	if max <= min {
		time.Sleep(time.Duration(min) * time.Millisecond)
		return
	}
	duration := cryptoRandInt(min, max)
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

// LongRandomSleep for major action transitions (page loads, etc)
func LongRandomSleep(minSec, maxSec int) {
	RandomSleep(minSec*1000, maxSec*1000)
}

// WaitForElement waits for an element with timeout, returns nil if not found
// This is safer than MustElement which panics
func WaitForElement(page *rod.Page, selector string, timeout time.Duration) (*rod.Element, error) {
	page = page.Timeout(timeout)
	el, err := page.Element(selector)
	if err != nil {
		return nil, err
	}
	return el, nil
}

// WaitForElementByText waits for an element containing specific text
func WaitForElementByText(page *rod.Page, elementType, text string, timeout time.Duration) (*rod.Element, error) {
	page = page.Timeout(timeout)
	el, err := page.ElementR(elementType, text)
	if err != nil {
		return nil, err
	}
	return el, nil
}

// SafeClick attempts to click with error handling and retries
func SafeClick(page *rod.Page, el *rod.Element, maxRetries int) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		err := HumanClick(page, el)
		if err == nil {
			return nil
		}
		lastErr = err

		// Wait before retry
		RandomSleep(500, 1000)
	}

	return lastErr
}

// IsElementVisible checks if an element is visible on page
func IsElementVisible(el *rod.Element) bool {
	if el == nil {
		return false
	}

	visible, err := el.Visible()
	if err != nil {
		return false
	}

	return visible
}
