package actor

import (
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
)

func BuildActors(parsedConfig *config.Config) map[string]actor.Actor {
	handler := map[string]actor.Actor{}

	buildRest(handler, parsedConfig.Actors.Rest)
	buildCombined(handler, parsedConfig.Actors.Combined)
	buildGfxBlink(handler, parsedConfig.Actors.GfxBlink)

	return handler
}
