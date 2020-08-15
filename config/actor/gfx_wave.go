package actor

import (
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
	"time"
)

func buildGfxWave(actors map[string]actor.Actor, gfxActors map[string]config.GfxWaveActor) {
	for actorName, gActor := range gfxActors {
		handler := &actor.GfxWave{
			Square: gActor.Square,
			Color:  colorOrDefault(gActor.Color, pad.ColorHighGreen),
			Delay:  gActor.Delay,
		}
		if handler.Delay == 0 {
			handler.Delay = 500 * time.Millisecond
		}

		actors[actorName] = handler
	}
}
