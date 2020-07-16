package gfx

import (
	"context"
	"go.uber.org/zap"
	"time"
)

var activeSequenceCancelFunc context.CancelFunc

// Sequence draw each given frame after the given delay.
func (e Renderer) Sequence(delay time.Duration, frames ...[]FramePixel) context.CancelFunc {
	ctx, cancelFunc := context.WithCancel(context.Background())

	//check if there is already an active sequence, if so cancel them
	if activeSequenceCancelFunc != nil {
		zap.L().Debug("Stop active sequence.")
		activeSequenceCancelFunc()
	}
	activeSequenceCancelFunc = cancelFunc

	draw := func() {
		if err := e.Clear(); err != nil {
			zap.L().Error("Could not clear pad!", zap.Error(err))
		}
		if err := e.Pattern(frames[0]...); err != nil {
			zap.L().Error("Could not draw pattern!", zap.Error(err))
		}

		//remove first one, so the next one is the first on next tick
		frames = frames[1:]
	}

	ticker := time.NewTicker(delay)
	go func() {
		defer ticker.Stop()

		//draw the first frame immediately
		draw()

		for {
			select {
			case <-ticker.C:
				if len(frames) > 0 {
					draw()
				} else {
					return
				}
			case <-ctx.Done():
				return
			}
		}

	}()

	return cancelFunc
}
