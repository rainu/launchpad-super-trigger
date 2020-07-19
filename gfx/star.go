package gfx

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/pad"
	"time"
)

// Star will animate a star which begins at the given position
func (e Renderer) Star(x, y int, color pad.Color, delay time.Duration) context.CancelFunc {
	spreadSequence := buildStarSeq(x, y, color)

	return e.Sequence(delay, spreadSequence...)
}

func buildStarSeq(x int, y int, color pad.Color) Sequence {
	sequences := make([]Sequence, 8)

	sequences[0] = buildBeamSeq(x, y, color, NorthSpreadDirection)
	sequences[1] = buildBeamSeq(x, y, color, NorthWestSpreadDirection)
	sequences[2] = buildBeamSeq(x, y, color, NorthEastSpreadDirection)
	sequences[3] = buildBeamSeq(x, y, color, EastSpreadDirection)
	sequences[4] = buildBeamSeq(x, y, color, SouthSpreadDirection)
	sequences[5] = buildBeamSeq(x, y, color, SouthEastSpreadDirection)
	sequences[6] = buildBeamSeq(x, y, color, SouthWestSpreadDirection)
	sequences[7] = buildBeamSeq(x, y, color, WestSpreadDirection)

	maxFrameCount := 0
	for _, sequence := range sequences {
		if maxFrameCount < len(sequence) {
			maxFrameCount = len(sequence)
		}
	}

	seq := make(Sequence, 0, maxFrameCount)

	for frame := 0; frame < maxFrameCount; frame++ {
		joinedFrame := make(Frame, 0, maxFrameCount)

		for _, subSeq := range sequences {
			if len(subSeq) > frame {
				joinedFrame = append(joinedFrame, subSeq[frame]...)
			} else {
				//duplicate the last frame of sequence
				joinedFrame = append(joinedFrame, subSeq[len(subSeq)-1]...)
			}
		}

		joinedFrame = mergePixels(joinedFrame...)
		seq = append(seq, joinedFrame)
	}

	return seq
}
