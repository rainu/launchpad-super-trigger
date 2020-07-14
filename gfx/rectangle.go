package gfx

import (
	"errors"
	"github.com/rainu/launchpad-super-trigger/pad"
)

// Rectangle draws a rectangle with the given line and fill color
func (e Renderer) Rectangle(x0, y0, x1, y1 int, fill, line pad.Color) error {
	if err := e.Fill(x0, y0, x1, y1, line); err != nil {
		return err
	}

	return e.Fill(x0+1, y0+1, x1-1, y1-1, fill)
}

// RectangleQuadrant draws a rectangle in the given quadrant with the given color settings
//  2.| 1.
// ---+---
//  3.| 4.
func (e Renderer) RectangleQuadrant(q Quadrant, fill, line pad.Color) error {
	switch q {
	case FirstQuadrant:
		return e.Rectangle(4, 0, 7, 3, fill, line)
	case SecondQuadrant:
		return e.Rectangle(0, 0, 3, 3, fill, line)
	case ThirdQuadrant:
		return e.Rectangle(0, 4, 3, 7, fill, line)
	case ForthQuadrant:
		return e.Rectangle(4, 4, 7, 7, fill, line)
	default:
		return errors.New("invalid quadrant")
	}
}
