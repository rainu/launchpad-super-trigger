package sensor

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/config/connection/mqtt"
	"github.com/rainu/launchpad-super-trigger/sensor"
)

func buildMqttSensors(sensors map[string]Sensor, generalSettings config.General, mqttSensors map[string]config.MQTTSensor, mqttConnections map[string]mqtt.Client) {
	for sensorName, mqttSensor := range mqttSensors {
		client := mqttConnections[mqttSensor.Connection]
		s := &sensor.MQTT{
			Client:       client,
			Topic:        mqttSensor.Topic,
			QOS:          mqttSensor.QOS,
			MessageStore: generateStore(generalSettings, sensorName),
		}

		sensors[sensorName] = Sensor{
			Sensor:     s,
			Extractors: buildExtractors(mqttSensor.DataPoints),
		}

		//If mqtt connection were lost, we have to resubscribe topics after
		//reconnection! so here we register a lister which will do so!
		client.AddListener(&connectionListener{
			sensor: s,
		})
	}
}

type connectionListener struct {
	sensor *sensor.MQTT
}

func (c *connectionListener) OnConnect(client MQTT.Client) {
	c.sensor.Reinitialise()
}

func (c *connectionListener) OnConnectionLost(client MQTT.Client, err error) {
	c.sensor.Purge()
}
