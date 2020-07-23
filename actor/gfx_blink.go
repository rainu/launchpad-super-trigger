package actor

import (
	"github.com/rainu/launchpad-super-trigger/gfx"
	"github.com/rainu/launchpad-super-trigger/pad"
	"time"
)

type GfxBlink struct {
	ColorOn  pad.Color
	ColorOff pad.Color

	Interval time.Duration
	Duration time.Duration
}

func (g *GfxBlink) Do(ctx Context) error {
	renderer := gfx.Renderer{ctx.Lighter}

	return renderer.BlinkBlocking(ctx.HitX, ctx.HitY, g.ColorOn, g.ColorOff, g.Interval, g.Duration, ctx.Context)
}
