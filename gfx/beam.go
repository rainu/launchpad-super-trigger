package gfx

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/pad"
	"time"
)

// Beam will animate a beam which begins at the given position and moves to the given direction
func (e Renderer) Beam(x, y int, dir SpreadDirection, color pad.Color, delay time.Duration) context.CancelFunc {
	spreadSequence := buildBeamSeq(x, y, color, dir)

	return e.Sequence(delay, spreadSequence...)
}

func buildBeamSeq(x, y int, color pad.Color, dir SpreadDirection) [][]FramePixel {
	spreadSequence := buildSpreadSeq(x, y, color, dir)

	//we can reuse spread sequence .. the different is that we have to join the frame
	//with their previous ones

	for i := len(spreadSequence) - 1; i > 0; i-- {
		frameCount := 0
		for j := i; j > 0; j-- {
			frameCount += len(spreadSequence[i-j])
		}

		joinFrame := make([]FramePixel, 0, frameCount)
		for j := i; j > 0; j-- {
			joinFrame = append(joinFrame, spreadSequence[i-j]...)
		}

		spreadSequence[i] = joinPixels(joinFrame...)
	}

	return spreadSequence
}
