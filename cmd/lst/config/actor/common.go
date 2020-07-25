package actor

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
)

func BuildActors(parsedConfig *config.Config, mqttConnections map[string]MQTT.Client) map[string]actor.Actor {
	handler := map[string]actor.Actor{}

	buildRest(handler, parsedConfig.Actors.Rest)
	buildMqtt(handler, parsedConfig.Actors.Mqtt, mqttConnections)
	buildCombined(handler, parsedConfig.Actors.Combined)
	buildGfxBlink(handler, parsedConfig.Actors.GfxBlink)
	buildGfxWave(handler, parsedConfig.Actors.GfxWave)

	return handler
}
