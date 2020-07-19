package gfx

import (
	"errors"
	"github.com/rainu/launchpad-super-trigger/pad"
)

// Fill fills the given rectangle with the given color
func (e Renderer) Fill(x0, y0, x1, y1 int, color pad.Color) error {
	frame := buildFill(x0, y0, x1, y1, color)
	for _, pixel := range frame {
		if err := pixel.Light(e); err != nil {
			return err
		}
	}

	return nil
}

func buildFill(x0, y0, x1, y1 int, color pad.Color) Frame {
	frame := make(Frame, 0, 64)

	appendPixel := func(x, y int) {
		if x < minX || y < minY || x > maxX || y > maxY {
			return
		}

		frame = append(frame, FramePixel{
			X:     x,
			Y:     y,
			Color: color,
		})
	}

	if x1 >= x0 && y1 >= y0 {
		// +-->
		// |
		// v
		for x := x0; x <= x1; x++ {
			for y := y0; y <= y1; y++ {
				appendPixel(x, y)
			}
		}
	} else if x1 <= x0 && y1 >= y0 {
		// <--+
		//    |
		//    v
		for x := x0; x >= x1; x-- {
			for y := y0; y <= y1; y++ {
				appendPixel(x, y)
			}
		}
	} else if x1 >= x0 && y1 <= y0 {
		// ^
		// |
		// +-->
		for x := x0; x <= x1; x++ {
			for y := y0; y >= y1; y-- {
				appendPixel(x, y)
			}
		}
	} else {
		//    ^
		//    |
		// <--+
		for x := x1; x >= x0; x-- {
			for y := y0; y >= y1; y-- {
				appendPixel(x, y)
			}
		}
	}

	return frame
}

// FillQuadrant fills the given quadrant with the given color
//  2.| 1.
// ---+---
//  3.| 4.
func (e Renderer) FillQuadrant(q Quadrant, color pad.Color) error {
	switch q {
	case FirstQuadrant:
		return e.Fill(4, minY, maxY, 3, color)
	case SecondQuadrant:
		return e.Fill(minX, minY, 3, 3, color)
	case ThirdQuadrant:
		return e.Fill(minX, 4, 3, maxY, color)
	case ForthQuadrant:
		return e.Fill(4, 4, maxX, maxY, color)
	default:
		return errors.New("invalid quadrant")
	}
}

// FillHorizontalLine fills the at the given position a line of the given color with the given length
func (e Renderer) FillHorizontalLine(x, y, length int, color pad.Color) error {
	return e.Fill(x, y, x+length-1, y, color)
}

// FillVerticalLine fills the at the given position a line of the given color with the given length
func (e Renderer) FillVerticalLine(x, y, length int, color pad.Color) error {
	return e.Fill(x, y, x, y+length-1, color)
}
