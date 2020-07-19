package gfx

import (
	"github.com/rainu/launchpad-super-trigger/pad"
)

// Circle draws a circle with the given line and fill color
func (e Renderer) Circle(x, y, r int, fill, line pad.Color) error {
	circle := buildFillCircle(x, y, r, fill, line)
	return e.Pattern(circle...)
}

func buildFillCircle(x, y, r int, fill, line pad.Color) Frame {
	cLine := buildCircle(x, y, r, line, true)
	cFill := make([]Frame, 0, padHeight)

	i := y - r + 1
	if i < minY {
		i = minY
	}
	minX := map[int]int{}
	maxX := map[int]int{}

	for ; i < y+r && i <= maxY; i++ {
		for _, pixel := range cLine {
			if pixel.Y == i {
				if min, found := minX[i]; !found || min > pixel.X {
					minX[i] = pixel.X
				}
				if max, found := maxX[i]; !found || max < pixel.X {
					maxX[i] = pixel.X
				}
			}
		}
	}
	for i, _ := range minX {
		if minX[i]+1 <= maxX[i]-1 {
			cFill = append(cFill, buildFill(minX[i]+1, i, maxX[i]-1, i, fill))
		}
	}

	return overrideFrames(cLine, overrideFrames(cFill...))
}

func buildCircle(x, y, r int, line pad.Color, includeDead bool) Frame {
	if r <= 0 {
		//only one point
		return Frame{{
			X:     x,
			Y:     y,
			Color: line,
		}}
	}

	frame := make(Frame, 0, 64)

	appendPixel := func(x, y int) {
		if includeDead {
			frame = append(frame, FramePixel{x, y, line})
			return
		}

		if x >= minX && x <= maxX && y >= minY && y <= maxY {
			frame = append(frame, FramePixel{x, y, line})
		}
	}

	// Bresenham algorithm - https://de.wikipedia.org/wiki/Bresenham-Algorithmus
	x1, y1, err := -r, 0, 2-2*r
	for {
		appendPixel(x-x1, y+y1)
		appendPixel(x-y1, y-x1)
		appendPixel(x+x1, y-y1)
		appendPixel(x+y1, y+x1)

		r = err
		if r > x1 {
			x1++
			err += x1*2 + 1
		}
		if r <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}

	return frame
}
