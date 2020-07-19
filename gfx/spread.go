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

func buildSpreadSeq(x int, y int, color pad.Color, spreadDir SpreadDirection) Sequence {
	seq := make(Sequence, 0, 9)

	seq = append(seq, Frame{{x, y, color}})

	switch spreadDir {
	case NorthSpreadDirection:
		for i := y - 1; i >= minY; i-- {
			frame := Frame{{x, i, color}}
			seq = append(seq, frame)
		}
	case SouthSpreadDirection:
		for i := y + 1; i < maxY; i++ {
			frame := Frame{{x, i, color}}
			seq = append(seq, frame)
		}
	case EastSpreadDirection:
		for i := x + 1; i < maxX; i++ {
			frame := Frame{{i, y, color}}
			seq = append(seq, frame)
		}
	case WestSpreadDirection:
		for i := x - 1; i >= minX; i-- {
			frame := Frame{{i, y, color}}
			seq = append(seq, frame)
		}
	case NorthEastSpreadDirection:
		i := y - 1
		j := x + 1
		for i >= minY && j <= maxX {
			frame := Frame{{j, i, color}}
			seq = append(seq, frame)

			i--
			j++
		}
	case SouthEastSpreadDirection:
		i := y + 1
		j := x + 1
		for i <= maxY && j <= maxX {
			frame := Frame{{j, i, color}}
			seq = append(seq, frame)

			i++
			j++
		}
	case SouthWestSpreadDirection:
		i := y + 1
		j := x - 1
		for i <= maxY && j >= minX {
			frame := Frame{{j, i, color}}
			seq = append(seq, frame)

			i++
			j--
		}
	case NorthWestSpreadDirection:
		i := y - 1
		j := x - 1
		for i >= minY && j >= minX {
			frame := Frame{{j, i, color}}
			seq = append(seq, frame)

			i--
			j--
		}
	}

	lastFrame := seq[len(seq)-1]
	seq = append(seq, Frame{{lastFrame[0].X, lastFrame[0].Y, pad.ColorOff}})

	return seq
}
