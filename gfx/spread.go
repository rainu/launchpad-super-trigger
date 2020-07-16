package gfx

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/pad"
	"time"
)

type SpreadDirection int

const (
	NorthSpreadDirection     SpreadDirection = 1
	NorthWestSpreadDirection SpreadDirection = 2
	NorthEastSpreadDirection SpreadDirection = 3
	EastSpreadDirection      SpreadDirection = 4
	SouthSpreadDirection     SpreadDirection = 5
	SouthEastSpreadDirection SpreadDirection = 6
	SouthWestSpreadDirection SpreadDirection = 7
	WestSpreadDirection      SpreadDirection = 8
)

// Spread will start at given point and move into the give direction with the specific delay
func (e Renderer) Spread(x, y int, dir SpreadDirection, color pad.Color, delay time.Duration) context.CancelFunc {
	spreadSequence := buildSpreadSeq(x, y, color, dir)

	return e.Sequence(delay, spreadSequence...)
}

func buildSpreadSeq(x int, y int, color pad.Color, spreadDir SpreadDirection) [][]FramePixel {
	seq := make([][]FramePixel, 0, 9)

	seq = append(seq, []FramePixel{{x, y, color}})

	switch spreadDir {
	case NorthSpreadDirection:
		for i := y - 1; i >= 0; i-- {
			frame := []FramePixel{{x, i, color}}
			seq = append(seq, frame)
		}
	case SouthSpreadDirection:
		for i := y + 1; i < 8; i++ {
			frame := []FramePixel{{x, i, color}}
			seq = append(seq, frame)
		}
	case EastSpreadDirection:
		for i := x + 1; i < 8; i++ {
			frame := []FramePixel{{i, y, color}}
			seq = append(seq, frame)
		}
	case WestSpreadDirection:
		for i := x - 1; i >= 0; i-- {
			frame := []FramePixel{{i, y, color}}
			seq = append(seq, frame)
		}
	case NorthEastSpreadDirection:
		i := y - 1
		j := x + 1
		for i >= 0 && j < 8 {
			frame := []FramePixel{{j, i, color}}
			seq = append(seq, frame)

			i--
			j++
		}
	case SouthEastSpreadDirection:
		i := y + 1
		j := x + 1
		for i < 8 && j < 8 {
			frame := []FramePixel{{j, i, color}}
			seq = append(seq, frame)

			i++
			j++
		}
	case SouthWestSpreadDirection:
		i := y + 1
		j := x - 1
		for i < 8 && j >= 0 {
			frame := []FramePixel{{j, i, color}}
			seq = append(seq, frame)

			i++
			j--
		}
	case NorthWestSpreadDirection:
		i := y - 1
		j := x - 1
		for i >= 0 && j >= 0 {
			frame := []FramePixel{{j, i, color}}
			seq = append(seq, frame)

			i--
			j--
		}
	}

	lastFrame := seq[len(seq)-1]
	seq = append(seq, []FramePixel{{lastFrame[0].X, lastFrame[0].Y, pad.ColorOff}})

	return seq
}
