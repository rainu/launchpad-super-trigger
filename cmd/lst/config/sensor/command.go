package sensor

import (
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/sensor"
)

func buildCommandSensors(sensors map[string]Sensor, commandSensors map[string]config.CommandSensor) {
	for sensorName, commandSensor := range commandSensors {
		s := &sensor.Command{
			Name:      commandSensor.Name,
			Arguments: commandSensor.Arguments,
			Interval:  commandSensor.Interval,
		}

		sensors[sensorName] = Sensor{
			Sensor:     s,
			Extractors: buildExtractors(commandSensor.DataPoints),
		}
	}
}
