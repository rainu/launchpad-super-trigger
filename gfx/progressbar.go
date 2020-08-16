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

	length := xUntil - xFrom
	p := (length * percent) / 100
	pixel := make([]FramePixel, 0, length)

	if dir == AscDirection {
		for x := xFrom; x < xUntil; x++ {
			pixel = append(pixel, FramePixel{X: x, Y: y, Color: empty})
		}
	} else {
		for x := xUntil - 1; x >= xFrom; x-- {
			pixel = append(pixel, FramePixel{X: x, Y: y, Color: empty})
		}
	}

	for i := 0; i < p && i < len(pixel); i++ {
		pixel[i].Color = fill
	}

	return e.Pattern(pixel...)
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

	length := yUntil - yFrom
	p := (length * percent) / 100
	pixel := make([]FramePixel, 0, length)

	if dir == AscDirection {
		for y := yUntil - 1; y >= yFrom; y-- {
			pixel = append(pixel, FramePixel{X: x, Y: y, Color: empty})
		}
	} else {
		for y := yFrom; y < yUntil; y++ {
			pixel = append(pixel, FramePixel{X: x, Y: y, Color: empty})
		}
	}

	for i := 0; i < p && i < len(pixel); i++ {
		pixel[i].Color = fill
	}

	return e.Pattern(pixel...)
}
