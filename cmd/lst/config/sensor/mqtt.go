package sensor

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/sensor"
)

func BuildMqttSensors(mqttSensors map[string]config.MQTTSensor, mqttConnections map[string]MQTT.Client) map[string]Sensor {
	result := map[string]Sensor{}

	for sensorName, mqttSensor := range mqttSensors {
		result[sensorName] = Sensor{
			Sensor: &sensor.MQTT{
				Client: mqttConnections[mqttSensor.Connection],
				Topic:  mqttSensor.Topic,
				QOS:    mqttSensor.QOS,
			},
			Extractors: buildExtractors(mqttSensor.DataPoints),
		}
	}

	return result
}
