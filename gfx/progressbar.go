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
	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		percent = 100
	}

	if err := e.Fill(xFrom, y, xUntil-1, y, empty); err != nil {
		return err
	}
	if percent == 0 {
		return nil
	}

	length := xUntil - xFrom
	if length < 0 {
		length *= -1
	}

	var x0, x1 int
	if dir == AscDirection {
		x0 = xFrom
		p := (length * percent) / 100
		x1 = x0 + p

		if x1 > maxX {
			return nil
		}
	} else {
		x0 = xUntil - 1
		p := (length * percent) / 100
		x1 = x0 - p

		if x1 < minY {
			return nil
		}
	}

	return e.Fill(x0, y, x1, y, fill)
}

func (e Renderer) VerticalProgressbar(x, percent int, dir Direction, fill, empty pad.Color) error {
	return e.verticalQuadrantProgressbar(x, minY, padHeight, percent, dir, fill, empty)
}

func (e Renderer) VerticalQuadrantProgressbar(q Quadrant, x, percent int, dir Direction, fill, empty pad.Color) error {
	switch q {
	case FirstQuadrant:
		return e.verticalQuadrantProgressbar(x+4, minY, 4, percent, dir, fill, empty)
	case SecondQuadrant:
		return e.verticalQuadrantProgressbar(x, minY, 4, percent, dir, fill, empty)
	case ThirdQuadrant:
		return e.verticalQuadrantProgressbar(x, 4, padHeight, percent, dir, fill, empty)
	case ForthQuadrant:
		return e.verticalQuadrantProgressbar(x+4, 4, padHeight, percent, dir, fill, empty)
	default:
		return errors.New("invalid quadrant")
	}
}

func (e Renderer) verticalQuadrantProgressbar(x, yFrom, yUntil, percent int, dir Direction, fill, empty pad.Color) error {
	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		percent = 100
	}

	if err := e.Fill(x, yFrom, x, yUntil-1, empty); err != nil {
		return err
	}
	if percent == 0 {
		return nil
	}
	length := yUntil - yFrom
	if length < 0 {
		length *= -1
	}

	var y0, y1 int
	if dir == AscDirection {
		y0 = yUntil - 1
		p := (length * percent) / 100
		y1 = y0 - p

		if y1 < minY {
			return nil
		}
	} else {
		y0 = yFrom
		p := (length * percent) / 100
		y1 = yFrom + p

		if y1 > maxY {
			return nil
		}
	}

	return e.Fill(x, y0, x, y1, fill)
}
