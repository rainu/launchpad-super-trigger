package meta

import (
	"github.com/rainu/launchpad-super-trigger/actor"
)

type LaunchpadSuperTriggerNavigationModeSwitch struct {
	Mode byte
}

func (l *LaunchpadSuperTriggerNavigationModeSwitch) Do(ctx actor.Context) error {
	return ctx.LST.SwitchNavigationMode(l.Mode)
}
