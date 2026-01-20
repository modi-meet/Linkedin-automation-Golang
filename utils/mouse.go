package utils

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

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

// ============================================================
// TECHNIQUE 1: HUMAN-LIKE MOUSE MOVEMENT (Bézier + overshoot)
// ============================================================

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

	jitter := 6.0
	targetX += (cryptoRandFloat() * jitter) - (jitter / 2)
	targetY += (cryptoRandFloat() * jitter) - (jitter / 2)

	moveMouseWithBezier(page, targetX, targetY)

	if cryptoRandInt(0, 10) < 3 {
		overshootX := targetX + (cryptoRandFloat()*20 - 10)
		overshootY := targetY + (cryptoRandFloat()*20 - 10)
		page.Mouse.MoveTo(proto.NewPoint(overshootX, overshootY))
		RandomSleep(30, 80)
		page.Mouse.MoveTo(proto.NewPoint(targetX, targetY))
	}

	RandomSleep(60, 150)

	if err := page.Mouse.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	RandomSleep(100, 300)
	return nil
}

func moveMouseWithBezier(page *rod.Page, targetX, targetY float64) {
	startX := cryptoRandFloat()*200 + 100
	startY := cryptoRandFloat()*200 + 100

	ctrlX1 := startX + (targetX-startX)*0.25 + (cryptoRandFloat()*80 - 40)
	ctrlY1 := startY + (targetY-startY)*0.25 + (cryptoRandFloat()*80 - 40)
	ctrlX2 := startX + (targetX-startX)*0.75 + (cryptoRandFloat()*50 - 25)
	ctrlY2 := startY + (targetY-startY)*0.75 + (cryptoRandFloat()*50 - 25)

	steps := 30 + cryptoRandInt(0, 20)
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		t = easeInOutQuad(t)

		x := cubicBezier(startX, ctrlX1, ctrlX2, targetX, t)
		y := cubicBezier(startY, ctrlY1, ctrlY2, targetY, t)

		if cryptoRandInt(0, 100) < 5 {
			x += cryptoRandFloat()*4 - 2
			y += cryptoRandFloat()*4 - 2
		}

		page.Mouse.MoveTo(proto.NewPoint(x, y))

		baseDelay := 4 + cryptoRandInt(0, 6)
		if t < 0.15 || t > 0.85 {
			baseDelay += cryptoRandInt(3, 8)
		}
		time.Sleep(time.Duration(baseDelay) * time.Millisecond)
	}

	page.Mouse.MoveTo(proto.NewPoint(targetX, targetY))
}

func easeInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - math.Pow(-2*t+2, 2)/2
}

func cubicBezier(p0, p1, p2, p3, t float64) float64 {
	return math.Pow(1-t, 3)*p0 +
		3*math.Pow(1-t, 2)*t*p1 +
		3*(1-t)*math.Pow(t, 2)*p2 +
		math.Pow(t, 3)*p3
}

// ============================================================
// TECHNIQUE 2: RANDOMIZED TIMING PATTERNS
// ============================================================

func RandomSleep(min, max int) {
	if max <= min {
		time.Sleep(time.Duration(min) * time.Millisecond)
		return
	}

	baseDuration := cryptoRandInt(min, max)

	if cryptoRandInt(0, 100) < 15 {
		baseDuration += cryptoRandInt(100, 500)
	}

	time.Sleep(time.Duration(baseDuration) * time.Millisecond)
}

func LongRandomSleep(minSec, maxSec int) {
	baseSleep := cryptoRandInt(minSec*1000, maxSec*1000)

	if cryptoRandInt(0, 100) < 20 {
		baseSleep += cryptoRandInt(500, 2000)
	}

	time.Sleep(time.Duration(baseSleep) * time.Millisecond)
}

func ThinkingPause() {
	RandomSleep(800, 2500)
}

// ============================================================
// TECHNIQUE 3: RANDOM SCROLLING BEHAVIOR
// ============================================================

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
		chunk := cryptoRandInt(30, 150)
		if scrolled+chunk > totalDistance {
			chunk = totalDistance - scrolled
		}

		speed := 0.5 + cryptoRandFloat()*1.5
		actualChunk := int(float64(chunk) * speed)
		if actualChunk < 10 {
			actualChunk = 10
		}

		page.Mouse.Scroll(0, float64(actualChunk*direction), 1)
		scrolled += actualChunk

		RandomSleep(50, 200)

		if cryptoRandInt(0, 100) < 8 {
			backScroll := cryptoRandInt(20, 80)
			page.Mouse.Scroll(0, float64(-backScroll*direction), 1)
			RandomSleep(100, 300)
			page.Mouse.Scroll(0, float64(backScroll*direction), 1)
			scrolled += 0
		}

		if cryptoRandInt(0, 100) < 12 {
			RandomSleep(400, 1200)
		}
	}
}

func ScrollToElement(page *rod.Page, el *rod.Element) error {
	if el == nil {
		return nil
	}

	RandomSleep(200, 500)

	shape, err := el.Shape()
	if err != nil {
		return err
	}

	box := shape.Box()
	viewportHeight := 1080.0
	scrollAmount := box.Y - (viewportHeight / 3)

	if scrollAmount > 50 {
		HumanScroll(page, int(scrollAmount))
	}

	RandomSleep(300, 600)
	return nil
}

// ============================================================
// TECHNIQUE 4: REALISTIC TYPING WITH TYPOS & CORRECTIONS
// ============================================================

func HumanType(page *rod.Page, el *rod.Element, text string) error {
	if el == nil || text == "" {
		return nil
	}

	if err := el.Focus(); err != nil {
		return err
	}

	RandomSleep(200, 500)

	chars := []rune(text)
	i := 0

	for i < len(chars) {
		if cryptoRandInt(0, 100) < 3 && i > 0 && i < len(chars)-1 {
			wrongChar := getAdjacentKey(chars[i])
			el.MustInput(string(wrongChar))
			RandomSleep(80, 200)

			ThinkingPause()

			page.Keyboard.MustType('\b')
			RandomSleep(50, 150)
		}

		el.MustInput(string(chars[i]))

		delay := getTypingDelay(chars[i])
		RandomSleep(delay, delay+80)

		if cryptoRandInt(0, 100) < 5 {
			RandomSleep(300, 800)
		}

		i++
	}

	RandomSleep(200, 500)
	return nil
}

func getTypingDelay(char rune) int {
	switch {
	case char == ' ':
		return cryptoRandInt(80, 180)
	case char == '.' || char == ',' || char == '!' || char == '?':
		return cryptoRandInt(120, 250)
	case char >= 'A' && char <= 'Z':
		return cryptoRandInt(100, 200)
	default:
		return cryptoRandInt(50, 130)
	}
}

func getAdjacentKey(char rune) rune {
	adjacentKeys := map[rune][]rune{
		'a': {'s', 'q', 'z'},
		'b': {'v', 'n', 'g'},
		'c': {'x', 'v', 'd'},
		'd': {'s', 'f', 'e'},
		'e': {'w', 'r', 'd'},
		'f': {'d', 'g', 'r'},
		'g': {'f', 'h', 't'},
		'h': {'g', 'j', 'y'},
		'i': {'u', 'o', 'k'},
		'j': {'h', 'k', 'u'},
		'k': {'j', 'l', 'i'},
		'l': {'k', 'o', 'p'},
		'm': {'n', 'k', 'j'},
		'n': {'b', 'm', 'h'},
		'o': {'i', 'p', 'l'},
		'p': {'o', 'l'},
		'q': {'w', 'a'},
		'r': {'e', 't', 'f'},
		's': {'a', 'd', 'w'},
		't': {'r', 'y', 'g'},
		'u': {'y', 'i', 'j'},
		'v': {'c', 'b', 'f'},
		'w': {'q', 'e', 's'},
		'x': {'z', 'c', 's'},
		'y': {'t', 'u', 'h'},
		'z': {'a', 'x'},
	}

	lowerChar := char
	if char >= 'A' && char <= 'Z' {
		lowerChar = char + 32
	}

	if adjacent, ok := adjacentKeys[lowerChar]; ok {
		return adjacent[cryptoRandInt(0, len(adjacent))]
	}

	return char
}

// ============================================================
// TECHNIQUE 5: MOUSE HOVERING & WANDERING
// ============================================================

func RandomHover(page *rod.Page) {
	x := cryptoRandFloat()*800 + 100
	y := cryptoRandFloat()*600 + 100

	moveMouseWithBezier(page, x, y)
	RandomSleep(200, 600)
}

func HoverOverElement(page *rod.Page, el *rod.Element) error {
	if el == nil {
		return nil
	}

	shape, err := el.Shape()
	if err != nil {
		return err
	}

	box := shape.Box()
	targetX := box.X + (box.Width / 2) + (cryptoRandFloat()*10 - 5)
	targetY := box.Y + (box.Height / 2) + (cryptoRandFloat()*10 - 5)

	moveMouseWithBezier(page, targetX, targetY)

	RandomSleep(300, 800)
	return nil
}

func IdleMouseMovement(page *rod.Page) {
	movements := cryptoRandInt(2, 5)
	for i := 0; i < movements; i++ {
		x := cryptoRandFloat()*400 + 200
		y := cryptoRandFloat()*300 + 150

		steps := cryptoRandInt(5, 15)
		for j := 0; j < steps; j++ {
			x += cryptoRandFloat()*20 - 10
			y += cryptoRandFloat()*20 - 10
			page.Mouse.MoveTo(proto.NewPoint(x, y))
			time.Sleep(time.Duration(cryptoRandInt(30, 80)) * time.Millisecond)
		}

		RandomSleep(500, 1500)
	}
}

// ============================================================
// TECHNIQUE 6: HELPER FUNCTIONS
// ============================================================

func WaitForElement(page *rod.Page, selector string, timeout time.Duration) (*rod.Element, error) {
	return page.Timeout(timeout).Element(selector)
}

func WaitForElementByText(page *rod.Page, elementType, text string, timeout time.Duration) (*rod.Element, error) {
	return page.Timeout(timeout).ElementR(elementType, text)
}

func SafeClick(page *rod.Page, el *rod.Element, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := HumanClick(page, el); err == nil {
			return nil
		} else {
			lastErr = err
		}
		RandomSleep(500, 1000)
	}
	return lastErr
}

func IsElementVisible(el *rod.Element) bool {
	if el == nil {
		return false
	}
	visible, err := el.Visible()
	return err == nil && visible
}

func SimulateReading(page *rod.Page) {
	readTime := cryptoRandInt(2000, 5000)

	movements := readTime / 1000
	for i := 0; i < movements; i++ {
		if cryptoRandInt(0, 100) < 40 {
			HumanScroll(page, cryptoRandInt(100, 300))
		}
		RandomSleep(800, 1500)
	}
}
