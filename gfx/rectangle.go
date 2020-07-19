package gfx

import (
	"errors"
	"github.com/rainu/launchpad-super-trigger/pad"
)

// Rectangle draws a rectangle with the given line and fill color
func (e Renderer) Rectangle(x0, y0, x1, y1 int, fill, line pad.Color) error {
	rect := buildRectangle(x0, y0, x1, y1, fill, line)
	return e.Pattern(rect...)
}

func buildRectangle(x0, y0, x1, y1 int, fill, line pad.Color) Frame {
	f0 := buildFill(x0, y0, x1, y1, line)
	f1 := buildFill(x0+1, y0+1, x1-1, y1-1, fill)
	rect := overrideFrames(f0, f1)

	return rect
}

// RectangleQuadrant draws a rectangle in the given quadrant with the given color settings
//  2.| 1.
// ---+---
//  3.| 4.
func (e Renderer) RectangleQuadrant(q Quadrant, fill, line pad.Color) error {
	switch q {
	case FirstQuadrant:
		return e.Rectangle(4, minY, maxX, 3, fill, line)
	case SecondQuadrant:
		return e.Rectangle(minX, minY, 3, 3, fill, line)
	case ThirdQuadrant:
		return e.Rectangle(minX, 4, 3, maxY, fill, line)
	case ForthQuadrant:
		return e.Rectangle(4, 4, maxX, maxY, fill, line)
	default:
		return errors.New("invalid quadrant")
	}
}
