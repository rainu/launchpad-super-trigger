package plotter

import "github.com/rainu/launchpad-super-trigger/pad"

type StaticExpression struct {
	ActivationColor pad.Color
	Matches         func([]byte) bool
}

type Static struct {
	X            int
	Y            int
	Expressions  []StaticExpression
	DefaultColor *pad.Color
}

func (s Static) Plot(ctx Context) error {
	for _, expression := range s.Expressions {
		if expression.Matches(ctx.Data) {

			//found an expression which matches
			return expression.ActivationColor.Light(ctx.Lighter, s.X, s.Y)
		}
	}

	if s.DefaultColor != nil {
		return s.DefaultColor.Light(ctx.Lighter, s.X, s.Y)
	}

	return nil
}
