package actor

import (
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
)

func buildCommand(actors map[string]actor.Actor, commandActors map[string]config.CommandActor) {
	for actorName, cmdActor := range commandActors {
		handler := &actor.Command{
			Name:          cmdActor.Name,
			Arguments:     cmdActor.Arguments,
			AppendContext: cmdActor.AppendContext,
		}

		actors[actorName] = handler
	}
}
