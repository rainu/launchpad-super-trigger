package gfx

import (
	"errors"
	"github.com/rainu/launchpad-super-trigger/pad"
)

func (e Renderer) HorizontalProgressbar(y, percent int, dir Direction, fill, empty pad.Color) error {
	return e.horizontalQuadrantProgressbar(y, minX, padLength, percent, dir, fill, empty)
}

func (e Renderer) HorizontalQuadrantProgressbar(q Quadrant, y, percent int, dir Direction, fill, empty pad.Color) error {
	switch q {
	case FirstQuadrant:
		return e.horizontalQuadrantProgressbar(y, 4, padLength, percent, dir, fill, empty)
	case SecondQuadrant:
		return e.horizontalQuadrantProgressbar(y, minX, 4, percent, dir, fill, empty)
	case ThirdQuadrant:
		return e.horizontalQuadrantProgressbar(y+4, minX, 4, percent, dir, fill, empty)
	case ForthQuadrant:
		return e.horizontalQuadrantProgressbar(y+4, 4, padLength, percent, dir, fill, empty)
	default:
		return errors.New("invalid quadrant")
	}
}

func (e Renderer) horizontalQuadrantProgressbar(y, xFrom, xUntil, percent int, dir Direction, fill, empty pad.Color) error {
	if err := e.Fill(xFrom, y, xUntil-1, y, empty); err != nil {
		return err
	}
	length := xUntil - xFrom

	x0 := xFrom
	x1 := xFrom + ((length * percent) / 100) - 1

	if dir == DescDirection {
		x0 = xUntil - 1
		x1 = xUntil - ((length * percent) / 100)
	}

	return e.Fill(x0, y, x1, y, fill)
}

func (e Renderer) VerticalProgressbar(y, percent int, dir Direction, fill, empty pad.Color) error {
	return e.verticalQuadrantProgressbar(y, padHeight, minY, percent, dir, fill, empty)
}

func (e Renderer) VerticalQuadrantProgressbar(q Quadrant, x, percent int, dir Direction, fill, empty pad.Color) error {
	switch q {
	case FirstQuadrant:
		return e.verticalQuadrantProgressbar(x+4, 4, minY, percent, dir, fill, empty)
	case SecondQuadrant:
		return e.verticalQuadrantProgressbar(x, 4, minY, percent, dir, fill, empty)
	case ThirdQuadrant:
		return e.verticalQuadrantProgressbar(x, padHeight, 4, percent, dir, fill, empty)
	case ForthQuadrant:
		return e.verticalQuadrantProgressbar(x+4, padHeight, 4, percent, dir, fill, empty)
	default:
		return errors.New("invalid quadrant")
	}
}

func (e Renderer) verticalQuadrantProgressbar(x, yFrom, yUntil, percent int, dir Direction, fill, empty pad.Color) error {
	if err := e.Fill(x, yFrom-1, x, yUntil, empty); err != nil {
		return err
	}
	length := yUntil - yFrom

	y0 := yFrom - 1
	y1 := yFrom + ((length * percent) / 100)

	if dir == DescDirection {
		y0 = yUntil
		y1 = yUntil - ((length * percent) / 100) - 1
	}

	return e.Fill(x, y0, x, y1, fill)
}
