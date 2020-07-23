package config

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
	"go.uber.org/zap"
	"sync"
)

type coord struct {
	X int
	Y int
}

type ColorSettings struct {
	Ready    pad.Color
	Progress pad.Color
	Success  pad.Color
	Failed   pad.Color
}

var defaultColorSettings = ColorSettings{
	Ready:    pad.ColorHighGreen,
	Progress: pad.ColorHighOrange,
	Success:  pad.ColorHighGreen,
	Failed:   pad.ColorHighRed,
}

type pageHandler struct {
	pageNumber    pad.PageNumber
	page          config.Page
	delegates     map[coord]actor.Actor
	colorSettings map[coord]ColorSettings

	activeProcess      map[coord]context.CancelFunc
	activeProcessMutex sync.RWMutex
}

func (p *pageHandler) Init(actors map[string]actor.Actor) {
	p.delegates = map[coord]actor.Actor{}
	p.colorSettings = map[coord]ColorSettings{}
	p.activeProcess = map[coord]context.CancelFunc{}

	for coord, trigger := range p.page.Trigger {
		c := convertCoordinate(coord)

		if trigger.ColorSettings != nil {
			p.colorSettings[c] = convertColorSettings(trigger.ColorSettings)
		} else {
			p.colorSettings[c] = defaultColorSettings
		}

		p.delegates[c] = actors[trigger.Actor]
	}
}

func (p *pageHandler) OnTrigger(lighter pad.Lighter, number pad.PageNumber, x int, y int) error {
	p.activeProcessMutex.RLock()
	_, found := p.activeProcess[coord{x, y}]
	p.activeProcessMutex.RUnlock()

	if found {
		zap.L().Warn("Ignore trigger because this action is already running")
		return nil
	}

	if delegate, found := p.delegates[coord{x, y}]; found {
		go func(p *pageHandler, delegate actor.Actor) {
			ctx, cancelFunc := context.WithCancel(context.Background())
			p.activeProcessMutex.Lock()
			p.activeProcess[coord{x, y}] = cancelFunc
			p.activeProcessMutex.Unlock()

			if err := p.colorSettings[coord{x, y}].Progress.Light(lighter, x, y); err != nil {
				zap.L().Debug("Could not light the pad!", zap.Error(err))
			}
			err := delegate.Do(actor.Context{
				Lighter: lighter,
				Context: ctx,
				Page:    number,
				HitX:    x,
				HitY:    y,
			})
			if err != nil {
				zap.L().Error("Actor returns an error: %w", zap.Error(err))

				if err := p.colorSettings[coord{x, y}].Failed.Light(lighter, x, y); err != nil {
					zap.L().Debug("Could not light the pad!", zap.Error(err))
				}
			} else {
				if err := p.colorSettings[coord{x, y}].Success.Light(lighter, x, y); err != nil {
					zap.L().Debug("Could not light the pad!", zap.Error(err))
				}
			}

			p.activeProcessMutex.Lock()
			delete(p.activeProcess, coord{x, y})
			p.activeProcessMutex.Unlock()
		}(p, delegate)
	}

	return nil
}

func (p *pageHandler) OnPageEnter(lighter pad.Lighter, number pad.PageNumber) error {
	for c, settings := range p.colorSettings {
		if err := settings.Ready.Light(lighter, c.X, c.Y); err != nil {
			zap.L().Debug("Could not light the pad!", zap.Error(err))
		}
	}

	return nil
}

func (p *pageHandler) OnPageLeave(lighter pad.Lighter, number pad.PageNumber) error {
	//on page leave close all running processes
	p.activeProcessMutex.RLock()
	for _, cancelFunc := range p.activeProcess {
		cancelFunc()
	}
	p.activeProcess = map[coord]context.CancelFunc{}
	p.activeProcessMutex.RUnlock()

	return nil
}

func convertCoordinate(coordinate config.Coordinate) coord {
	x, y, err := coordinate.Coordinate()
	if err != nil {
		panic(err)
	}

	return coord{x, y}
}

func convertColorSettings(settings *config.ColorSettings) ColorSettings {
	return ColorSettings{
		Ready:    convertColor(settings.Ready),
		Progress: convertColor(settings.Progress),
		Success:  convertColor(settings.Success),
		Failed:   convertColor(settings.Failed),
	}
}

func convertColor(color config.Color) pad.Color {
	c, err := color.Color()
	if err != nil {
		panic(err)
	}
	return c
}
