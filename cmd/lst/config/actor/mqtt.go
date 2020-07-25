package actor

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
)

func buildMqtt(actors map[string]actor.Actor, mqttActors map[string]config.MQTTActor, mqttConnections map[string]MQTT.Client) {
	for actorName, mqttActor := range mqttActors {
		handler := &actor.Mqtt{
			Client:   mqttConnections[mqttActor.Connection],
			Topic:    mqttActor.Topic,
			QOS:      mqttActor.QOS,
			Retained: false,
		}

		if mqttActor.BodyRaw != "" {
			handler.Body = rawBody(mqttActor.BodyRaw)
		} else if mqttActor.BodyB64 != "" {
			handler.Body = b64Body(mqttActor.BodyB64)
		} else if mqttActor.BodyPath != "" {
			handler.Body = fileBody(mqttActor.BodyPath)
		}

		actors[actorName] = handler
	}
}
