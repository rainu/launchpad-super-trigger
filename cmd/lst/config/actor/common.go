package actor

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/cmd/lst/config/sensor"
	"github.com/rainu/launchpad-super-trigger/config"
)

func BuildActors(parsedConfig *config.Config, sensors map[string]sensor.Sensor, mqttConnections map[string]MQTT.Client) map[string]actor.Actor {
	handler := map[string]actor.Actor{}

	buildRest(handler, parsedConfig.Actors.Rest)
	buildMqtt(handler, parsedConfig.Actors.Mqtt, mqttConnections)
	buildCommand(handler, parsedConfig.Actors.Command)
	buildCombined(handler, parsedConfig.Actors.Combined)
	buildGfxBlink(handler, parsedConfig.Actors.GfxBlink)
	buildGfxWave(handler, parsedConfig.Actors.GfxWave)

	//must be last
	buildConditional(handler, sensors, parsedConfig.Actors.Conditional)

	return handler
}
