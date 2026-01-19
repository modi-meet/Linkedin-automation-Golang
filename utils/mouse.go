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

func HumanType(page *rod.Page, el *rod.Element, text string) error {
	if el == nil {
		return nil
	}

	if err := el.Focus(); err != nil {
		return err
	}

	RandomSleep(200, 400)

	for _, char := range text {
		if err := el.Input(string(char)); err != nil {
			return err
		}

		if char == ' ' || char == '.' || char == ',' || char == '!' || char == '?' {
			RandomSleep(100, 300)
		} else {
			RandomSleep(50, 180)
		}
	}

	RandomSleep(200, 500)
	return nil
}

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

	jitter := 8.0
	targetX += (cryptoRandFloat() * jitter) - (jitter / 2)
	targetY += (cryptoRandFloat() * jitter) - (jitter / 2)

	moveMouseWithBezier(page, targetX, targetY)
	RandomSleep(80, 200)

	if err := page.Mouse.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	RandomSleep(150, 400)
	return nil
}

func moveMouseWithBezier(page *rod.Page, targetX, targetY float64) {
	startX := cryptoRandFloat() * 100
	startY := cryptoRandFloat() * 100

	ctrlX1 := startX + (targetX-startX)*0.3 + (cryptoRandFloat()*100 - 50)
	ctrlY1 := startY + (targetY-startY)*0.3 + (cryptoRandFloat()*100 - 50)
	ctrlX2 := startX + (targetX-startX)*0.7 + (cryptoRandFloat()*60 - 30)
	ctrlY2 := startY + (targetY-startY)*0.7 + (cryptoRandFloat()*60 - 30)

	steps := 25 + cryptoRandInt(0, 15)

	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := cubicBezier(startX, ctrlX1, ctrlX2, targetX, t)
		y := cubicBezier(startY, ctrlY1, ctrlY2, targetY, t)

		page.Mouse.MoveTo(proto.NewPoint(x, y))

		baseDelay := 5 + cryptoRandInt(0, 8)
		if t < 0.2 || t > 0.8 {
			baseDelay += cryptoRandInt(2, 6)
		}
		time.Sleep(time.Duration(baseDelay) * time.Millisecond)
	}

	page.Mouse.MoveTo(proto.NewPoint(targetX, targetY))
}

func cubicBezier(p0, p1, p2, p3, t float64) float64 {
	return math.Pow(1-t, 3)*p0 +
		3*math.Pow(1-t, 2)*t*p1 +
		3*(1-t)*math.Pow(t, 2)*p2 +
		math.Pow(t, 3)*p3
}

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
		chunk := cryptoRandInt(50, 200)
		if scrolled+chunk > totalDistance {
			chunk = totalDistance - scrolled
		}

		page.Mouse.Scroll(0, float64(chunk*direction), 1)
		scrolled += chunk

		RandomSleep(100, 400)

		if cryptoRandInt(0, 10) < 2 {
			RandomSleep(500, 1500)
		}
	}
}

func RandomSleep(min, max int) {
	if max <= min {
		time.Sleep(time.Duration(min) * time.Millisecond)
		return
	}
	duration := cryptoRandInt(min, max)
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

func LongRandomSleep(minSec, maxSec int) {
	RandomSleep(minSec*1000, maxSec*1000)
}

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
