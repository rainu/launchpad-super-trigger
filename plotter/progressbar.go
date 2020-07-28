package plotter

import (
	"fmt"
	"github.com/rainu/launchpad-super-trigger/gfx"
	"github.com/rainu/launchpad-super-trigger/pad"
	"strconv"
)

type Progressbar struct {
	X         int
	Y         int
	Min       float64
	Max       float64
	Vertical  bool
	Quadrant  gfx.Quadrant
	Direction gfx.Direction
	Fill      pad.Color
	Empty     pad.Color
}

func (p Progressbar) Plot(ctx Context) error {
	renderer := gfx.Renderer{ctx.Lighter}

	parsedFloat, err := strconv.ParseFloat(string(ctx.Data), 64)
	if err != nil {
		return fmt.Errorf("can not parse as integer value: %w", err)
	}

	//((input - min) * 100) / (max - min)
	percentage := ((parsedFloat - p.Min) * 100) / (p.Max - p.Min)

	if p.Vertical {
		if p.Quadrant > 0 {
			return renderer.VerticalQuadrantProgressbar(p.Quadrant, p.X, int(percentage), p.Direction, p.Fill, p.Empty)
		} else {
			return renderer.VerticalProgressbar(p.X, int(percentage), p.Direction, p.Fill, p.Empty)
		}
	} else {
		if p.Quadrant > 0 {
			return renderer.HorizontalQuadrantProgressbar(p.Quadrant, p.Y, int(percentage), p.Direction, p.Fill, p.Empty)
		} else {
			return renderer.HorizontalProgressbar(p.Y, int(percentage), p.Direction, p.Fill, p.Empty)
		}
	}
}
