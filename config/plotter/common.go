package plotter

import (
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/plotter"
)

func BuildPlotter(plotters config.Plotters) map[plotter.Plotter]config.Datapoint {
	result := map[plotter.Plotter]config.Datapoint{}

	buildProgressbar(result, plotters.Progressbar)
	buildStatic(result, plotters.Static)

	return result
}
