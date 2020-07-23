package gfx

import (
	"context"
	"fmt"
	"github.com/rainu/launchpad-super-trigger/pad"
	"go.uber.org/zap"
	"time"
)

var activeBlinks = map[coord]context.CancelFunc{}

func (e Renderer) Blink(x, y int, on, off pad.Color, interval, duration time.Duration) context.CancelFunc {
	ctx, cancelFunc := context.WithCancel(context.Background())

	//check if there is already an active blink, if so cancel them
	if activeCancelFunc, exists := activeBlinks[coord{x, y}]; exists {
		zap.L().Debug(fmt.Sprintf("Stop active blink at %d, %d", x, y))
		activeCancelFunc()
	}
	activeBlinks[coord{x, y}] = cancelFunc

	go func() {
		err := e.BlinkBlocking(x, y, on, off, interval, duration, ctx)
		if err != nil {
			zap.L().Error("Error while gfx blink!", zap.Error(err))
		}
	}()

	return cancelFunc
}

func (e Renderer) BlinkBlocking(x, y int, on, off pad.Color, interval, duration time.Duration, ctx context.Context) error {
	ticker := time.NewTicker(interval)
	timer := time.NewTimer(duration)

	defer ticker.Stop()
	defer timer.Stop()
	defer off.Light(e, x, y)

	currentColor := on

	for {
		select {
		case <-ticker.C:
			if err := currentColor.Light(e, x, y); err != nil {
				return fmt.Errorf("could not light the pad: %w", err)
			}

			if currentColor.Ordinal() == on.Ordinal() {
				currentColor = off
			} else {
				currentColor = on
			}
		case <-timer.C:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
