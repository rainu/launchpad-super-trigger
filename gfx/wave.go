package gfx

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/pad"
	"time"
)

// WaveSquare will animate a rectangle wave which begin at given point
func (e Renderer) WaveSquare(x, y int, color pad.Color, delay time.Duration) context.CancelFunc {
	seq := make(Sequence, 0, 9)

	firstEmpty := true
	for i := 0; i < 9; i++ {
		rect := buildRectangle(x-i, y-i, x+i, y+i, pad.ColorOff, color)

		if !rect.HasOnlyColor(pad.ColorOff) {
			seq = append(seq, rect)
		} else if firstEmpty {
			seq = append(seq, rect)
			firstEmpty = false
		}
	}

	return e.Sequence(delay, seq...)
}

// WaveCircle will animate a circle wave which begin at given point
func (e Renderer) WaveCircle(x, y int, color pad.Color, delay time.Duration) context.CancelFunc {
	seq := make(Sequence, 0, 9)

	firstEmpty := true
	for i := 0; ; i++ {
		circle := buildCircle(x, y, i, color, false)

		if !circle.HasOnlyColor(pad.ColorOff) {
			seq = append(seq, circle)
		} else if firstEmpty {
			seq = append(seq, circle)
			firstEmpty = false
		} else {
			break
		}
	}

	return e.Sequence(delay, seq...)
}
