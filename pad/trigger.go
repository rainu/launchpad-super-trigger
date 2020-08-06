package pad

import (
	"context"
	"fmt"
	"github.com/rainu/launchpad"
	"gitlab.com/gomidi/midi"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

// TriggerHandleFunc will called each time a hit was made.
// If the hit was outside the trigger area x and y will be 255!
type TriggerHandleFunc func(lighter Lighter, page PageNumber, x, y int) error

type LaunchpadSuperTrigger struct {
	pad         launchpad.Launchpad
	driver      midi.Driver
	lighter     Lighter
	currentPage *Page
	specials    *special
	handle      TriggerHandleFunc
}

func NewLaunchpadSuperTrigger(driver midi.Driver, handler TriggerHandleFunc) (*LaunchpadSuperTrigger, error) {
	pad, err := launchpad.NewLaunchpad(driver)
	if err != nil {
		return nil, err
	}

	page := NewPage(0)
	special := newSpecial()

	return &LaunchpadSuperTrigger{
		pad:    pad,
		driver: driver,
		lighter: &triggerAreaLighter{
			page:    page,
			special: special,
			delegate: &threadSafeLighter{
				mux:      sync.Mutex{},
				delegate: pad,
			},
		},
		currentPage: page,
		specials:    special,
		handle:      handler,
	}, nil
}

func (l *LaunchpadSuperTrigger) WaitForConnectionLost(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
outerLoop:
	for {
		select {
		case <-ticker.C:
			ins, err := l.driver.Ins()
			if err != nil {
				return
			}

			for i := range ins {
				if strings.Contains(ins[i].String(), "Launchpad S") {
					continue outerLoop
				}
			}
			return //no launchpad found -> connections lost
		case <-ctx.Done():
			return
		}
	}
}

func (l *LaunchpadSuperTrigger) Run(ctx context.Context) {
	if err := l.pad.Clear(); err != nil {
		zap.L().Error("Error while clearing the launchpad!", zap.Error(err))
		return
	}
	if err := l.handle(l.lighter, l.currentPage.Number(), 255, 255); err != nil {
		zap.L().Error("Error while handling hit!", zap.Error(err))
	}

	hitChannel, err := l.pad.ListenToHits()
	if err != nil {
		zap.L().Error("Error while initialise the launchpad hit listener!", zap.Error(err))
		return
	}

	for {
		select {
		case hit := <-hitChannel:
			zap.L().Debug(fmt.Sprintf("Incoming hit: { %d; %d }", hit.X, hit.Y))

			//ignore hit releases at the moment
			if !hit.Down {
				continue
			}

			if IsPageHit(hit) {
				l.applyPage(hit)

				zap.L().Info(fmt.Sprintf("Switched to page: %d", l.currentPage.Number()))
				if err := l.handle(l.lighter, l.currentPage.Number(), 255, 255); err != nil {
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
