package gfx

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/pad"
	"time"
)

// Boom will animate a explosion at given point with a given delay
func (e Renderer) Boom(x, y int, color pad.Color, delay time.Duration) context.CancelFunc {
	boomSequence := buildBoomSeq(x, y, color)

	return e.Sequence(delay, boomSequence...)
}

func buildBoomSeq(x int, y int, color pad.Color) Sequence {
	sequences := make([]Sequence, 8)

	sequences[0] = buildSpreadSeq(x, y, color, NorthSpreadDirection)
	sequences[1] = buildSpreadSeq(x, y, color, NorthWestSpreadDirection)
	sequences[2] = buildSpreadSeq(x, y, color, NorthEastSpreadDirection)
	sequences[3] = buildSpreadSeq(x, y, color, EastSpreadDirection)
	sequences[4] = buildSpreadSeq(x, y, color, SouthSpreadDirection)
	sequences[5] = buildSpreadSeq(x, y, color, SouthEastSpreadDirection)
	sequences[6] = buildSpreadSeq(x, y, color, SouthWestSpreadDirection)
	sequences[7] = buildSpreadSeq(x, y, color, WestSpreadDirection)

	seq := make(Sequence, 0, 9)

	for frame := 0; ; frame++ {
		joinedFrame := make(Frame, 0, 9)

		for _, subSeq := range sequences {
			if len(subSeq) > frame {
				joinedFrame = append(joinedFrame, subSeq[frame]...)
			}
		}

		if len(joinedFrame) == 0 {
			break
		}
		joinedFrame = mergePixels(joinedFrame...)
		seq = append(seq, joinedFrame)
	}

	return seq
}
