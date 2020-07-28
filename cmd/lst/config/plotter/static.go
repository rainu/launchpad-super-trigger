package plotter

import (
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
	"github.com/rainu/launchpad-super-trigger/plotter"
	"regexp"
	"strconv"
	"strings"
)

func buildStatic(allPlotter map[plotter.Plotter]string, staticPlotter []config.Static) {
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
				Matches:         buildEqExpressionFn(expression),
			})
		}
		for _, expression := range static.Expressions.Ne {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         buildNeExpressionFn(expression),
			})
		}
		for _, expression := range static.Expressions.Lt {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         buildLtExpressionFn(expression),
			})
		}
		for _, expression := range static.Expressions.Lte {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         buildLteExpressionFn(expression),
			})
		}
		for _, expression := range static.Expressions.Gt {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         buildGtExpressionFn(expression),
			})
		}
		for _, expression := range static.Expressions.Gte {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         buildGteExpressionFn(expression),
			})
		}
		for _, expression := range static.Expressions.Match {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         buildMatchExpressionFn(expression),
			})
		}
		for _, expression := range static.Expressions.NotMatch {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         buildNotMatchExpressionFn(expression),
			})
		}
		for _, expression := range static.Expressions.Contains {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         buildContainsExpressionFn(expression),
			})
		}
		for _, expression := range static.Expressions.NotContains {
			sb.Expressions = append(sb.Expressions, plotter.StaticExpression{
				ActivationColor: colorOrDefault(expression.ActivationColor, pad.ColorHighGreen),
				Matches:         buildNotContainsExpressionFn(expression),
			})
		}

		allPlotter[sb] = static.DataPoint
	}
}

func buildEqExpressionFn(expr config.StaticExpression) func([]byte) bool {
	isNumericValue := false
	fExpValue, err := strconv.ParseFloat(expr.Value, 64)
	if err == nil {
		isNumericValue = true
	}

	return func(value []byte) bool {
		if isNumericValue {
			fValue, err := strconv.ParseFloat(string(value), 64)
			if err != nil {
				return false
			}

			return fValue == fExpValue
		}

		return expr.Value == string(value)
	}
}

func buildNeExpressionFn(expr config.StaticExpression) func([]byte) bool {
	isNumericValue := false
	fExpValue, err := strconv.ParseFloat(expr.Value, 64)
	if err == nil {
		isNumericValue = true
	}

	return func(value []byte) bool {
		if isNumericValue {
			fValue, err := strconv.ParseFloat(string(value), 64)
			if err != nil {
				return false
			}

			return fValue != fExpValue
		}

		return expr.Value != string(value)
	}
}
func buildLtExpressionFn(expr config.StaticNumericExpression) func([]byte) bool {
	return func(value []byte) bool {
		fValue, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return fValue < expr.Value
		}

		return false
	}
}

func buildLteExpressionFn(expr config.StaticNumericExpression) func([]byte) bool {
	return func(value []byte) bool {
		fValue, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return fValue <= expr.Value
		}

		return false
	}
}

func buildGtExpressionFn(expr config.StaticNumericExpression) func([]byte) bool {
	return func(value []byte) bool {
		fValue, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return fValue > expr.Value
		}

		return false
	}
}

func buildGteExpressionFn(expr config.StaticNumericExpression) func([]byte) bool {
	return func(value []byte) bool {
		fValue, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return fValue >= expr.Value
		}

		return false
	}
}

func buildMatchExpressionFn(expr config.StaticMatchExpression) func([]byte) bool {
	re := regexp.MustCompile(expr.Value)
	return func(value []byte) bool {
		return re.Match(value)
	}
}

func buildNotMatchExpressionFn(expr config.StaticMatchExpression) func([]byte) bool {
	re := regexp.MustCompile(expr.Value)
	return func(value []byte) bool {
		return !re.Match(value)
	}
}

func buildContainsExpressionFn(expr config.StaticExpression) func([]byte) bool {
	return func(value []byte) bool {
		return strings.Contains(string(value), expr.Value)
	}
}

func buildNotContainsExpressionFn(expr config.StaticExpression) func([]byte) bool {
	return func(value []byte) bool {
		return !strings.Contains(string(value), expr.Value)
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
