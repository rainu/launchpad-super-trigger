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

	ticker := time.NewTicker(interval)
	timer := time.NewTimer(duration)
	go func() {
		defer ticker.Stop()
		defer timer.Stop()
		defer off.Light(e, x, y)

		currentColor := on

		for {
			select {
			case <-ticker.C:
				if err := currentColor.Light(e, x, y); err != nil {
					zap.L().Error("Could not light the pad!", zap.Error(err))
				}

				if currentColor.Ordinal() == on.Ordinal() {
					currentColor = off
				} else {
					currentColor = on
				}
			case <-timer.C:
				return
			case <-ctx.Done():
				return
			}
		}

	}()

	return cancelFunc
}
