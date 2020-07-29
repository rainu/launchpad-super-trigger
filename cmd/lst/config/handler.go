package config

import (
	"context"
	"fmt"
	"github.com/rainu/launchpad-super-trigger/actor"
	configSensor "github.com/rainu/launchpad-super-trigger/cmd/lst/config/sensor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
	"github.com/rainu/launchpad-super-trigger/plotter"
	"github.com/rainu/launchpad-super-trigger/sensor"
	"github.com/rainu/launchpad-super-trigger/sensor/data_extractor"
	"go.uber.org/zap"
	"strings"
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
	actors        map[coord]actor.Actor
	sensors       map[string]configSensor.Sensor
	plotters      map[plotter.Plotter]config.Datapoint
	colorSettings map[coord]ColorSettings

	activeProcess      map[coord]context.CancelFunc
	activeProcessMutex sync.RWMutex

	lastLighter pad.Lighter
}

func (p *pageHandler) Init(actors map[string]actor.Actor, sensors map[string]configSensor.Sensor, plotters map[plotter.Plotter]config.Datapoint) {
	p.actors = map[coord]actor.Actor{}
	p.sensors = sensors
	p.plotters = plotters
	p.colorSettings = map[coord]ColorSettings{}
	p.activeProcess = map[coord]context.CancelFunc{}

	for triggerCoords, trigger := range p.page.Trigger {
		coordinates := convertCoordinates(triggerCoords)
		for _, c := range coordinates {
			if trigger.ColorSettings != nil {
				p.colorSettings[c] = convertColorSettings(trigger.ColorSettings)
			} else {
				p.colorSettings[c] = defaultColorSettings
			}

			p.actors[c] = actors[trigger.Actor]
		}
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

	if delegate, found := p.actors[coord{x, y}]; found {
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
	p.lastLighter = lighter

	for c, settings := range p.colorSettings {
		if err := settings.Ready.Light(lighter, c.X, c.Y); err != nil {
			zap.L().Debug("Could not light the pad!", zap.Error(err))
		}
	}

	for _, s := range p.sensors {
		s.Sensor.AddCallback(p, p.OnData)
	}

	//plot available sensor data
	for _, s := range p.sensors {
		if len(s.Sensor.LastMessage()) > 0 {
			p.OnData(s.Sensor)
		}
	}

	return nil
}

func (p *pageHandler) OnPageLeave(lighter pad.Lighter, number pad.PageNumber) error {
	p.lastLighter = nil

	//on page leave close all running processes
	p.activeProcessMutex.RLock()
	for _, cancelFunc := range p.activeProcess {
		cancelFunc()
	}
	p.activeProcess = map[coord]context.CancelFunc{}
	p.activeProcessMutex.RUnlock()

	for _, s := range p.sensors {
		s.Sensor.RemoveCallback(p)
	}

	return nil
}

func (p *pageHandler) OnData(sensor sensor.Sensor) {
	if p.lastLighter == nil {
		return
	}

	sensorName := ""
	var extractors map[string]data_extractor.Extractor
	for name, s := range p.sensors {
		if s.Sensor == sensor {
			sensorName = name
			extractors = s.Extractors
			break
		}
	}

	for responsiblePlotter, dataPoint := range p.plotters {
		if strings.HasPrefix(dataPoint.Path(), sensorName+".") {
			dpName := dataPoint.Name()
			extract, err := extractors[dpName].Extract(sensor.LastMessage())

			if err != nil {
				zap.L().Warn(fmt.Sprintf("Could not extract datapoint '%s.%s': ", sensorName, dpName), zap.Error(err))
				continue
			}

			if responsiblePlotter == nil {
				zap.L().Fatal("Could not found corresponding plotter! This should never happen (validation failed?)")
			}

			err = responsiblePlotter.Plot(plotter.Context{
				Lighter: p.lastLighter,
				Page:    p.pageNumber,
				Data:    extract,
			})

			if err != nil {
				zap.L().Error("Error while plotting incoming sensor data point!", zap.Error(err))
			}
		}
	}
}

func convertCoordinate(coordinate config.Coordinate) coord {
	x, y, err := coordinate.Coordinate()
	if err != nil {
		panic(err)
	}

	return coord{x, y}
}

func convertCoordinates(coordinates config.Coordinates) []coord {
	cfgCoords, err := coordinates.Coordinates()
	if err != nil {
		panic(err)
	}

	result := make([]coord, 0, len(cfgCoords))
	for _, c := range cfgCoords {
		result = append(result, coord{c[0], c[1]})
	}

	return result
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
		return pad.ColorOff
	}
	return c
}
