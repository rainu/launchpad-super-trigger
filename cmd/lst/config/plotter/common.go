package plotter

import (
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/plotter"
)

func BuildPlotter(plotters config.Plotters) map[plotter.Plotter]string {
	result := map[plotter.Plotter]string{}

	buildProgressbar(result, plotters.Progressbar)

	return result
}
