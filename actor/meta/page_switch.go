package meta

import (
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/pad"
)

type LaunchpadSuperTriggerPageSwitch struct {
	Target pad.PageNumber
}

func (l *LaunchpadSuperTriggerPageSwitch) Do(ctx actor.Context) error {
	return ctx.LST.SwitchPage(l.Target)
}
