package pad

import (
	"context"
	"fmt"
	"github.com/rakyll/launchpad"
	"go.uber.org/zap"
)

type LaunchpadSuperTrigger struct {
	pad         *launchpad.Launchpad
	currentPage *Page
}

func NewLaunchpadSuperTrigger() (*LaunchpadSuperTrigger, error) {
	pad, err := launchpad.Open()
	if err != nil {
		return nil, err
	}

	return &LaunchpadSuperTrigger{
		pad:         pad,
		currentPage: NewPage(0),
	}, nil
}

func (l *LaunchpadSuperTrigger) Run(ctx context.Context) {
	l.pad.Clear()
	hitChannel := l.pad.Listen()

	for {
		select {
		case hit := <-hitChannel:
			zap.L().Debug(fmt.Sprintf("Incoming hit: { %d; %d }", hit.X, hit.Y))

			if IsPageHit(hit) {
				l.currentPage.Toggle(hit.X)
				l.currentPage.Light(l.pad)

				zap.L().Info(fmt.Sprintf("Switched to page: %d", l.currentPage.Number()))
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
