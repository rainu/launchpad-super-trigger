package actor

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/template"
	"go.uber.org/zap"
)

func buildMqtt(actors map[string]actor.Actor, mqttActors map[string]config.MQTTActor, mqttConnections map[string]MQTT.Client, engine *template.Engine) {
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
		} else if mqttActor.BodyTemplate != "" {
			if err := engine.RegisterTemplate(actorName, mqttActor.BodyTemplate); err != nil {
				zap.L().Fatal("Failed to parse template!", zap.Error(err))
			}

			handler.Body = templateBody(actorName, engine)
		}

		actors[actorName] = handler
	}
}
