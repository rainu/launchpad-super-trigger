package actor

import (
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"go.uber.org/zap"
)

type combinedActor interface {
	actor.Actor

	AddActor(actor actor.Actor)
}

func buildCombined(actors map[string]actor.Actor, combinedActors map[string]config.CombinedActor) {
	for actorName, cActor := range combinedActors {
		var handler combinedActor

		if cActor.Parallel {
			handler = &actor.ParallelActor{}
		} else {
			handler = &actor.SequentialActor{}
		}

		for _, actorName := range cActor.Actor {
			delegate, found := actors[actorName]
			if !found {
				zap.L().Fatal("No corresponding actor found: " + actorName)
			}

			handler.AddActor(delegate)
		}

		actors[actorName] = handler
	}
}
