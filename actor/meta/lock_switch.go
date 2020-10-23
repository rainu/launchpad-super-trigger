package meta

import (
	"github.com/rainu/launchpad-super-trigger/actor"
)

type LaunchpadSuperTriggerLockSwitch struct {
	Lock bool
}

func (l *LaunchpadSuperTriggerLockSwitch) Do(ctx actor.Context) error {
	return ctx.LST.SwitchLock(l.Lock)
}
