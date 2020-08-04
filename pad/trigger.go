package pad

import (
	"context"
	"fmt"
	"github.com/rakyll/launchpad"
	"go.uber.org/zap"
	"sync"
)

// TriggerHandleFunc will called each time a hit was made.
// If the hit was outside the trigger area x and y will be -1!
type TriggerHandleFunc func(lighter Lighter, page PageNumber, x int, y int) error

type LaunchpadSuperTrigger struct {
	pad         ExtendedLaunchpad
	lighter     Lighter
	currentPage *Page
	specials    *special
	handle      TriggerHandleFunc
}

func NewLaunchpadSuperTrigger(handler TriggerHandleFunc) (*LaunchpadSuperTrigger, error) {
	pad, err := launchpad.Open()
	if err != nil {
		return nil, err
	}
	extendedPad := newExtendedLaunchpad(pad)

	page := NewPage(0)
	special := newSpecial()

	return &LaunchpadSuperTrigger{
		pad: extendedPad,
		lighter: &triggerAreaLighter{
			page:    page,
			special: special,
			delegate: &threadSafeLighter{
				mux:      sync.Mutex{},
				delegate: extendedPad,
			},
		},
		currentPage: page,
		specials:    special,
		handle:      handler,
	}, nil
}

func (l *LaunchpadSuperTrigger) Run(ctx context.Context) {
	if err := l.pad.Clear(); err != nil {
		zap.L().Error("Error while clearing the launchpad!", zap.Error(err))
		return
	}
	if err := l.handle(l.lighter, l.currentPage.Number(), -1, -1); err != nil {
		zap.L().Error("Error while handling hit!", zap.Error(err))
	}

	hitChannel := l.pad.Listen()

	for {
		select {
		case hit := <-hitChannel:
			zap.L().Debug(fmt.Sprintf("Incoming hit: { %d; %d }", hit.X, hit.Y))

			if IsPageHit(hit) {
				l.applyPage(hit)

				zap.L().Info(fmt.Sprintf("Switched to page: %d", l.currentPage.Number()))
				if err := l.handle(l.lighter, l.currentPage.Number(), -1, -1); err != nil {
					zap.L().Error("Error while handling hit!", zap.Error(err))
				}
			} else if IsPadHit(hit) {
				if l.specials.locked {
					zap.L().Debug("Dont handle hit because pad is locked.")
					continue
				}

				if err := l.handle(l.lighter, l.currentPage.Number(), hit.X, hit.Y); err != nil {
					zap.L().Error("Error while handling hit!", zap.Error(err))
				}
			} else if IsMetaTextMarker(hit) {
				zap.L().Debug("Received a text end marker!")
			} else if IsSpecialVol(hit) {
				//change the navigationMode
				if err := l.specials.SwitchPageNavigationMode(l.pad); err != nil {
					zap.L().Error("Could not switch page navigation mode!", zap.Error(err))
				}
			} else if IsSpecialArm(hit) {
				if err := l.specials.ToggleLock(l.pad); err != nil {
					zap.L().Error("Could not switch arm mode!", zap.Error(err))
				}
			} else {
				continue
			}

			//light the page buttons AFTER the handler, so handlers could use "lighter.Clear()" without
			//fear to lose the page lights
			if err := l.currentPage.Light(l.lighter); err != nil {
				zap.L().Error("Could not light the pad!", zap.Error(err))
			}
		case <-ctx.Done():
			//context closed -> stop running
			return
		}
	}
}

func (l *LaunchpadSuperTrigger) applyPage(hit launchpad.Hit) {
	switch l.specials.pageNavigationMode {
	case pageNavigationBinary:
		l.currentPage.Toggle(hit.X)
	case pageNavigationToggle:
		l.currentPage.SetTo(hit.X)
	}
}

func (l *LaunchpadSuperTrigger) Close() error {
	return l.pad.Close()
}
