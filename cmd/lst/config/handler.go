package config

import (
	"context"
	"fmt"
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/actor/meta"
	"github.com/rainu/launchpad-super-trigger/config"
	configSensor "github.com/rainu/launchpad-super-trigger/config/sensor"
	"github.com/rainu/launchpad-super-trigger/pad"
	"github.com/rainu/launchpad-super-trigger/plotter"
	"github.com/rainu/launchpad-super-trigger/sensor"
	"github.com/rainu/launchpad-super-trigger/sensor/data_extractor"
	"go.uber.org/zap"
	"reflect"
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
	Ready:    pad.ColorGreen,
	Progress: pad.ColorOrange,
	Success:  pad.ColorGreen,
	Failed:   pad.ColorRed,
}

type pageHandler struct {
	pageNumber    pad.PageNumber
	page          config.Page
	actors        map[coord]actor.Actor
	sensors       map[string]configSensor.Sensor
	plotters      map[plotter.Plotter]config.Datapoint
	colorSettings map[coord]*ColorSettings

	activeProcess      map[coord]context.CancelFunc
	activeProcessMutex sync.RWMutex

	lastLighter pad.Lighter
}

func (p *pageHandler) Init(actors map[string]actor.Actor, sensors map[string]configSensor.Sensor, plotters map[plotter.Plotter]config.Datapoint) {
	p.actors = map[coord]actor.Actor{}
	p.sensors = sensors
	p.plotters = plotters
	p.colorSettings = map[coord]*ColorSettings{}
	p.activeProcess = map[coord]context.CancelFunc{}

	for triggerCoords, trigger := range p.page.Trigger {
		coordinates := convertCoordinates(triggerCoords)
		for _, c := range coordinates {
			if trigger.ColorSettings != nil {
				p.colorSettings[c] = convertColorSettings(trigger.ColorSettings)
			} else {
				p.colorSettings[c] = &defaultColorSettings
			}

			//special use case: if this actor is an page switcher
			hasPageSwitch := false
			if reflect.TypeOf(actors[trigger.Actor]) == reflect.TypeOf(&meta.LaunchpadSuperTriggerPageSwitch{}) {
				hasPageSwitch = true
			} else if reflect.TypeOf(actors[trigger.Actor]) == reflect.TypeOf(&actor.Sequential{}) ||
				reflect.TypeOf(actors[trigger.Actor]) == reflect.TypeOf(&actor.Parallel{}) ||
				reflect.TypeOf(actors[trigger.Actor]) == reflect.TypeOf(&actor.Conditional{}) {

				var metaActor actor.MetaActor
				var ok bool

				metaActor, ok = actors[trigger.Actor].(*actor.Sequential)
				if !ok {
					metaActor, ok = actors[trigger.Actor].(*actor.Parallel)
					if !ok {
						metaActor = actors[trigger.Actor].(*actor.Conditional)
					}
				}

				//check if any underlying actor is an page switcher
				hasPageSwitch = metaActor.HasActor(func(actor actor.Actor) bool {
					return reflect.TypeOf(actor) == reflect.TypeOf(meta.LaunchpadSuperTriggerPageSwitch{}) ||
						reflect.TypeOf(actor) == reflect.TypeOf(&meta.LaunchpadSuperTriggerPageSwitch{})
				})
			}

			if hasPageSwitch {
				p.colorSettings[c] = &ColorSettings{
					Ready:    p.colorSettings[c].Ready,
					Progress: p.colorSettings[c].Progress,
					Failed:   p.colorSettings[c].Failed,
					Success:  nil,
				}
			}

			p.actors[c] = actors[trigger.Actor]
		}
	}
}

func (p *pageHandler) OnTrigger(lighter pad.Lighter, lst *pad.LaunchpadSuperTrigger, number pad.PageNumber, x int, y int) error {
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

			colorEnabled := p.colorSettings[coord{x, y}] != nil

			if colorEnabled {
				if err := p.colorSettings[coord{x, y}].Progress.Light(lighter, x, y); err != nil {
					zap.L().Debug("Could not light the pad!", zap.Error(err))
				}
			}
			err := delegate.Do(actor.Context{
				Lighter: lighter,
				LST:     lst,
				Context: ctx,
				Page:    number,
				HitX:    x,
				HitY:    y,
			})
			if err != nil {
				zap.L().Error("Actor returns an error: %w", zap.Error(err))

				if colorEnabled {
					if err := p.colorSettings[coord{x, y}].Failed.Light(lighter, x, y); err != nil {
						zap.L().Debug("Could not light the pad!", zap.Error(err))
					}
				}
			} else {
				if colorEnabled {
					successColor := p.colorSettings[coord{x, y}].Success

					if successColor != nil {
						if err := successColor.Light(lighter, x, y); err != nil {
							zap.L().Debug("Could not light the pad!", zap.Error(err))
						}
					}
				}
			}

			p.activeProcessMutex.Lock()
			delete(p.activeProcess, coord{x, y})
			p.activeProcessMutex.Unlock()
		}(p, delegate)
	}

	return nil
}

func (p *pageHandler) OnPageEnter(lighter pad.Lighter, lst *pad.LaunchpadSuperTrigger, number pad.PageNumber) error {
	p.lastLighter = lighter

	for c, settings := range p.colorSettings {
		if settings == nil {
			continue
		}

		if err := settings.Ready.Light(lighter, c.X, c.Y); err != nil {
			zap.L().Debug("Could not light the pad!", zap.Error(err))
		}
	}

	for _, s := range p.sensors {
		s.Sensor.AddCallback(p, p.OnData)
	}

	//plot available sensor data
	for _, s := range p.sensors {
		p.OnData(s.Sensor)
	}

	return nil
}

func (p *pageHandler) OnPageLeave(lighter pad.Lighter, lst *pad.LaunchpadSuperTrigger, number pad.PageNumber) error {
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

			zap.L().Debug(fmt.Sprintf("Extraction of datapoint %s: %s", dataPoint.Path(), extract))

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

func convertColorSettings(settings *config.ColorSettings) *ColorSettings {
	if settings.Disable {
		return nil
	}

	return &ColorSettings{
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
