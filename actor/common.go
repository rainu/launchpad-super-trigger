package actor

import (
	"context"
	"github.com/rainu/launchpad-super-trigger/pad"
)

type Context struct {
	Lighter pad.Lighter
	Context context.Context
	Page    pad.PageNumber
	HitX    int
	HitY    int
}

type Actor interface {
	Do(ctx Context) error
}
