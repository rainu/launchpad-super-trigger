package gfx

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type Sequence []Frame

var activeSequenceCancelFunc context.CancelFunc

// Sequence draw each given frame after the given delay.
func (e Renderer) Sequence(delay time.Duration, frames ...Frame) context.CancelFunc {
	ctx, cancelFunc := context.WithCancel(context.Background())

	//check if there is already an active sequence, if so cancel them
	if activeSequenceCancelFunc != nil {
		zap.L().Debug("Stop active sequence.")
		activeSequenceCancelFunc()
	}
	activeSequenceCancelFunc = cancelFunc

	go func() {
		err := e.SequenceBlocking(delay, ctx, frames...)
		if err != nil {
			zap.L().Error("Error while gfx sequence!", zap.Error(err))
		}
	}()

	return cancelFunc
}

func (e Renderer) SequenceBlocking(delay time.Duration, ctx context.Context, frames ...Frame) error {
	draw := func() error {
		if err := e.Clear(); err != nil {
			return fmt.Errorf("could not clear pad: %w", err)
		}
		if err := e.Pattern(frames[0]...); err != nil {
			return fmt.Errorf("could not draw pattern: %w", err)
		}

		//remove first one, so the next one is the first on next tick
		frames = frames[1:]
		return nil
	}

	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	//draw the first frame immediately
	if err := draw(); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if len(frames) > 0 {
				if err := draw(); err != nil {
					return err
				}
			} else {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
