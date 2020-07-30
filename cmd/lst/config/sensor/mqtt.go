package sensor

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/sensor"
)

func buildMqttSensors(sensors map[string]Sensor, mqttSensors map[string]config.MQTTSensor, mqttConnections map[string]MQTT.Client) {
	for sensorName, mqttSensor := range mqttSensors {
		sensors[sensorName] = Sensor{
			Sensor: &sensor.MQTT{
				Client: mqttConnections[mqttSensor.Connection],
				Topic:  mqttSensor.Topic,
				QOS:    mqttSensor.QOS,
			},
			Extractors: buildExtractors(mqttSensor.DataPoints),
		}
	}
}
