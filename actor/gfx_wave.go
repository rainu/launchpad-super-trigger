package actor

import (
	"github.com/rainu/launchpad-super-trigger/gfx"
	"github.com/rainu/launchpad-super-trigger/pad"
	"time"
)

type GfxWave struct {
	Square bool
	Color  pad.Color
	Delay  time.Duration
}

func (g *GfxWave) Do(ctx Context) error {
	renderer := gfx.Renderer{ctx.Lighter}

	if g.Square {
		return renderer.WaveSquareBlocking(ctx.HitX, ctx.HitY, g.Color, g.Delay, ctx.Context)
	}
	return renderer.WaveCircleBlocking(ctx.HitX, ctx.HitY, g.Color, g.Delay, ctx.Context)
}
