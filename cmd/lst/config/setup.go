package config

import (
	configActor "github.com/rainu/launchpad-super-trigger/cmd/lst/config/actor"
	"github.com/rainu/launchpad-super-trigger/cmd/lst/config/connection"
	configPlotter "github.com/rainu/launchpad-super-trigger/cmd/lst/config/plotter"
	configSensor "github.com/rainu/launchpad-super-trigger/cmd/lst/config/sensor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
	"github.com/rainu/launchpad-super-trigger/sensor"
	"io"
	"strings"
)

func ConfigureDispatcher(configReader io.Reader) (*pad.TriggerDispatcher, map[string]sensor.Sensor, error) {
	parsedConfig, err := config.ReadConfig(configReader)
	if err != nil {
		return nil, nil, err
	}

	dispatcher := &pad.TriggerDispatcher{}
	connections := connection.BuildMqttConnection(parsedConfig)
	actors := configActor.BuildActors(parsedConfig, connections)
	sensors := configSensor.BuildMqttSensors(parsedConfig.Sensors.Mqtt, connections)

	for pageNumber, page := range parsedConfig.Layout.Pages {
		handler := &pageHandler{
			pageNumber: pad.PageNumber(pageNumber),
			page:       page,
		}
		plotters := configPlotter.BuildPlotter(page.Plotter)

		usedSensors := map[string]configSensor.Sensor{}
		usedSensorNames := UsedSensors(page.Plotter)

		for _, sensorName := range usedSensorNames {
			usedSensors[sensorName] = sensors[sensorName]
		}

		handler.Init(actors, usedSensors, plotters)
		dispatcher.AddPageHandler(handler, handler.pageNumber)
	}

	deflatedSensors := map[string]sensor.Sensor{}
	for name, s := range sensors {
		deflatedSensors[name] = s.Sensor
	}

	return dispatcher, deflatedSensors, nil
}

func UsedSensors(p config.Plotters) []string {
	sensors := map[string]bool{}

	for _, progressbar := range p.Progressbar {
		sensorName := strings.Split(progressbar.DataPoint, ".")[0]
		sensors[sensorName] = true
	}

	result := make([]string, 0, len(sensors))
	for s, _ := range sensors {
		result = append(result, s)
	}

	return result
}
