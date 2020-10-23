package actor

import (
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/actor/meta"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
)

func buildPageSwitch(actors map[string]actor.Actor, psActors map[string]config.PageSwitchActor) {
	for actorName, psActor := range psActors {
		handler := &meta.LaunchpadSuperTriggerPageSwitch{
			Target: pad.PageNumber(psActor.Target.AsInt()),
		}

		actors[actorName] = handler
	}
}

func buildNavigationModeSwitch(actors map[string]actor.Actor, nmsActors map[string]config.NavigationModeSwitchActor) {
	for actorName, nmsActor := range nmsActors {
		handler := &meta.LaunchpadSuperTriggerNavigationModeSwitch{
			Mode: nmsActor.NavigationMode,
		}

		actors[actorName] = handler
	}
}

func buildLockSwitch(actors map[string]actor.Actor, lsActors map[string]config.LockSwitchActor) {
	for actorName, lsActor := range lsActors {
		handler := &meta.LaunchpadSuperTriggerLockSwitch{
			Lock: lsActor.Lock,
		}

		actors[actorName] = handler
	}
}
