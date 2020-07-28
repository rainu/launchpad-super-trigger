package plotter

import "github.com/rainu/launchpad-super-trigger/pad"

type Context struct {
	Lighter pad.Lighter
	Page    pad.PageNumber
	Data    []byte
}

type Plotter interface {
	Plot(ctx Context) error
}
