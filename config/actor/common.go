package actor

import (
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/config/connection/mqtt"
	"github.com/rainu/launchpad-super-trigger/config/sensor"
	"github.com/rainu/launchpad-super-trigger/template"
)

func BuildActors(parsedConfig *config.Config, sensors map[string]sensor.Sensor, templateEngine *template.Engine, mqttConnections map[string]mqtt.Client) map[string]actor.Actor {
	handler := map[string]actor.Actor{}

	buildRest(handler, parsedConfig.Actors.Rest, templateEngine)
	buildMqtt(handler, parsedConfig.Actors.Mqtt, mqttConnections, templateEngine)
	buildCommand(handler, parsedConfig.Actors.Command)
	buildCombined(handler, parsedConfig.Actors.Combined)
	buildPageSwitch(handler, parsedConfig.Actors.MetaPageSwitcher)
	buildNavigationModeSwitch(handler, parsedConfig.Actors.MetaNavigationModeSwitcher)
	buildLockSwitch(handler, parsedConfig.Actors.MetaLockerSwitcher)
	buildGfxBlink(handler, parsedConfig.Actors.GfxBlink)
	buildGfxWave(handler, parsedConfig.Actors.GfxWave)

	//must be last
	buildConditional(handler, sensors, parsedConfig.Actors.Conditional)

	return handler
}
