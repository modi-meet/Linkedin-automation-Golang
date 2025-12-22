package utils

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func HumanClick(page *rod.Page, el *rod.Element) {
	box := el.MustShape().Box()
	targetX := box.X + (box.Width / 2)
	targetY := box.Y + (box.Height / 2)

	jitter := 5.0
	targetX += (rand.Float64() * jitter) - (jitter / 2)
	targetY += (rand.Float64() * jitter) - (jitter / 2)

	moveMouseIdeally(page, targetX, targetY)

	RandomSleep(150, 300)
	page.Mouse.MustClick(proto.InputMouseButtonLeft)
}

// break a movement into small steps
func moveMouseIdeally(page *rod.Page, x, y float64) {
	const steps = 20

	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)

		startOffsetX := (rand.Float64()*30 - 15) * (1 - t)
		startOffsetY := (rand.Float64()*30 - 15) * (1 - t)

		page.Mouse.MoveTo(proto.NewPoint(x+startOffsetX, y+startOffsetY))
		time.Sleep(time.Duration(5+rand.Intn(10)) * time.Millisecond)
	}

	page.Mouse.MoveTo(proto.NewPoint(x, y))
}

func RandomSleep(min, max int) {
	if max <= min {
		time.Sleep(time.Duration(min) * time.Millisecond)
		return
	}
	duration := rand.Intn(max-min) + min
	time.Sleep(time.Duration(duration) * time.Millisecond)
}
