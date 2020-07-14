package pad

import (
	"context"
	"fmt"
	"github.com/rakyll/launchpad"
	"go.uber.org/zap"
	"sync"
)

type TriggerHandleFunc func(Lighter, PageNumber, int, int) error

type LaunchpadSuperTrigger struct {
	pad         *launchpad.Launchpad
	lighter     Lighter
	currentPage *Page
	handle      TriggerHandleFunc
}

func NewLaunchpadSuperTrigger(handler TriggerHandleFunc) (*LaunchpadSuperTrigger, error) {
	pad, err := launchpad.Open()
	if err != nil {
		return nil, err
	}

	return &LaunchpadSuperTrigger{
		pad: pad,
		lighter: &ThreadSafeLighter{
			mux:      sync.Mutex{},
			delegate: pad,
		},
		currentPage: NewPage(0),
		handle:      handler,
	}, nil
}

func (l *LaunchpadSuperTrigger) Run(ctx context.Context) {
	if err := l.pad.Clear(); err != nil {
		zap.L().Error("Error while clearing the launchpad!", zap.Error(err))
		return
	}

	hitChannel := l.pad.Listen()

	for {
		select {
		case hit := <-hitChannel:
			zap.L().Debug(fmt.Sprintf("Incoming hit: { %d; %d }", hit.X, hit.Y))

			if IsPageHit(hit) {
				l.currentPage.Toggle(hit.X)
				if err := l.pad.Clear(); err != nil {
					zap.L().Error("Could not clear the pad!", zap.Error(err))
				}
				if err := l.currentPage.Light(l.lighter); err != nil {
					zap.L().Error("Could not light the pad!", zap.Error(err))
				}

				zap.L().Info(fmt.Sprintf("Switched to page: %d", l.currentPage.Number()))
			} else if IsPadHit(hit) {
				if err := l.handle(l.lighter, l.currentPage.Number(), hit.X, hit.Y); err != nil {
					zap.L().Error("Error while handling hit!", zap.Error(err))
				}
			}
		case <-ctx.Done():
			//context closed -> stop running
			return
		}
	}
}

func (l *LaunchpadSuperTrigger) Close() error {
	return l.pad.Close()
}
