package config

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/config"
	configActor "github.com/rainu/launchpad-super-trigger/config/actor"
	connectionMqtt "github.com/rainu/launchpad-super-trigger/config/connection/mqtt"
	configPlotter "github.com/rainu/launchpad-super-trigger/config/plotter"
	configSensor "github.com/rainu/launchpad-super-trigger/config/sensor"
	"github.com/rainu/launchpad-super-trigger/pad"
	"github.com/rainu/launchpad-super-trigger/sensor"
	"github.com/rainu/launchpad-super-trigger/template"
	"go.uber.org/zap"
	"io"
	"reflect"
)

func ConfigureDispatcher(configReader ...io.Reader) (*pad.TriggerDispatcher, map[string]sensor.Sensor, config.General, error) {
	parsedConfig, err := config.ReadConfig(configReader...)
	if err != nil {
		return nil, nil, config.General{}, err
	}

	dispatcher := &pad.TriggerDispatcher{}
	connections := connectionMqtt.BuildMqttConnection(parsedConfig)
	sensors := configSensor.BuildSensors(parsedConfig.General, parsedConfig.Sensors, connections)
	templateEngine := setupTemplateEngine(sensors)
	actors := configActor.BuildActors(parsedConfig, sensors, templateEngine, connections)

	//establish mqtt connections after callbacks were registered
	for _, connection := range connections {
		if err := connection.Connect(context.Background()); err != nil {
			zap.L().Fatal("Error while connecting to mqtt broker: %s", zap.Error(err))
		}
	}

	for pageNumber, page := range parsedConfig.Layout.Pages {
		handler := &pageHandler{
			pageNumber: pad.PageNumber(pageNumber.AsInt()),
			page:       page,
		}
		plotters := configPlotter.BuildPlotter(page.Plotter)

		usedSensors := map[string]configSensor.Sensor{}
		usedSensorNames := getUsedSensorNames(page.Plotter)

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

	return dispatcher, deflatedSensors, parsedConfig.General, nil
}

func setupTemplateEngine(sensors map[string]configSensor.Sensor) *template.Engine {
	teSensors := map[string]template.Sensor{}

	for name, sensor := range sensors {
		teSensors[name] = template.Sensor{
			Sensor:     sensor.Sensor,
			Extractors: sensor.Extractors,
		}
	}

	return template.NewEngine(teSensors)
}

func getUsedSensorNames(p config.Plotters) []string {
	sensors := map[string]bool{}

	refPlotters := reflect.ValueOf(p)

	//Plotters
	for plottersField := 0; plottersField < refPlotters.NumField(); plottersField++ {
		if refPlotters.Field(plottersField).Kind() == reflect.Slice {

			//Plotters.Progressbar
			for i := 0; i < refPlotters.Field(plottersField).Len(); i++ {
				//Plotters.Progressbar[i]

				valPlotter := refPlotters.Field(plottersField).Index(i)

				for plotterField := 0; plotterField < valPlotter.NumField(); plotterField++ {
					if valPlotter.Field(plotterField).Type() == reflect.TypeOf(config.Datapoint("")) {
						dpPath := config.Datapoint(valPlotter.Field(plotterField).String())

						sensors[dpPath.Sensor()] = true
					}
				}
			}
		}
	}

	result := make([]string, 0, len(sensors))
	for s, _ := range sensors {
		result = append(result, s)
	}

	return result
}
