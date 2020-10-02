package plotter

import (
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/gfx"
	"github.com/rainu/launchpad-super-trigger/pad"
	"github.com/rainu/launchpad-super-trigger/plotter"
)

func buildProgressbar(allPlotter map[plotter.Plotter]config.Datapoint, progressPlotter []config.Progressbar) {
	for _, progressbar := range progressPlotter {
		pb := plotter.Progressbar{
			X:         progressbar.X,
			Y:         progressbar.Y,
			Quadrant:  progressbar.Quadrant,
			Min:       progressbar.Min,
			Direction: gfx.AscDirection,
			Vertical:  progressbar.Vertical,
			Fill:      colorOrDefault(progressbar.Fill, pad.ColorGreen),
			Empty:     colorOrDefault(progressbar.Empty, pad.ColorOff),
		}
		if progressbar.Max != nil {
			pb.Max = *progressbar.Max
		} else {
			pb.Max = 100
		}
		if progressbar.RightToLeft {
			pb.Direction = gfx.DescDirection
		}

		allPlotter[pb] = progressbar.DataPoint
	}
}

func colorOrDefault(color config.Color, defaultColor pad.Color) pad.Color {
	c, err := color.Color()
	if err != nil {
		return defaultColor
	}
	return c
}
