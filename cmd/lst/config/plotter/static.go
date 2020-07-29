package plotter

import (
	"github.com/rainu/launchpad-super-trigger/cmd/lst/config/expressions"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
	"github.com/rainu/launchpad-super-trigger/plotter"
)

func buildStatic(allPlotter map[plotter.Plotter]config.Datapoint, staticPlotter []config.Static) {
	for _, static := range staticPlotter {
		x, y, _ := static.Position.Coordinate()
		sb := &plotter.Static{
			X:            x,
			Y:            y,
			DefaultColor: pColorOrDefault(static.DefaultColor, nil),
			Expressions:  make([]plotter.StaticExpression, 0),
		}

		for _, expression := range static.Expressions.Eq {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildEqExpressionFn(expression.Value),
			})
		}
		for _, expression := range static.Expressions.Ne {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildNeExpressionFn(expression.Value),
			})
		}
		for _, expression := range static.Expressions.Lt {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildLtExpressionFn(expression.Value),
			})
		}
		for _, expression := range static.Expressions.Lte {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildLteExpressionFn(expression.Value),
			})
		}
		for _, expression := range static.Expressions.Gt {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildGtExpressionFn(expression.Value),
			})
		}
		for _, expression := range static.Expressions.Gte {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildGteExpressionFn(expression.Value),
			})
		}
		for _, expression := range static.Expressions.Match {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildMatchExpressionFn(expression.Value),
			})
		}
		for _, expression := range static.Expressions.NotMatch {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildNotMatchExpressionFn(expression.Value),
			})
		}
		for _, expression := range static.Expressions.Contains {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildContainsExpressionFn(expression.Value),
			})
		}
		for _, expression := range static.Expressions.NotContains {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         expressions.BuildNotContainsExpressionFn(expression.Value),
			})
		}

		allPlotter[sb] = static.DataPoint
	}
}

func pColorOrDefault(color *config.Color, defaultColor *pad.Color) *pad.Color {
	if color == nil {
		return defaultColor
	}

	c, err := color.Color()
	if err != nil {
		return defaultColor
	}
	return &c
}
