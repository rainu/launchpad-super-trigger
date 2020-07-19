package gfx

import "github.com/rainu/launchpad-super-trigger/pad"

type FramePixel struct {
	X     int
	Y     int
	Color pad.Color
}

func (f *FramePixel) Light(lighter pad.Lighter) error {
	return f.Color.Light(lighter, f.X, f.Y)
}

func (f *FramePixel) coord() coord {
	return coord{
		X: f.X,
		Y: f.Y,
	}
}

type Frame []FramePixel

func (f Frame) HasOnlyColor(c pad.Color) bool {
	for _, pixel := range f {
		if pixel.Color.Ordinal() != c.Ordinal() {
			return false
		}
	}

	return true
}

// Pattern will draw the given pixel
func (e Renderer) Pattern(frame ...FramePixel) error {
	for _, pixel := range frame {
		if err := pixel.Color.Light(e, pixel.X, pixel.Y); err != nil {
			return err
		}
	}

	return nil
}

func mergePixels(pixels ...FramePixel) Frame {
	alreadyKnownPositions := map[coord]bool{}
	joined := make(Frame, 0, len(pixels))

	for _, pixel := range pixels {
		if !alreadyKnownPositions[pixel.coord()] {
			alreadyKnownPositions[pixel.coord()] = true
			joined = append(joined, pixel)
		}
	}

	return joined
}

func overrideFrames(frames ...Frame) Frame {
	frameAsMap := map[coord]*FramePixel{}

	for _, frame := range frames {
		for i := range frame {
			pixel, exists := frameAsMap[frame[i].coord()]
			if !exists {
				frameAsMap[frame[i].coord()] = &FramePixel{
					X:     frame[i].X,
					Y:     frame[i].Y,
					Color: frame[i].Color,
				}
				continue
			}

			//override the existing color with this given color
			pixel.Color = frame[i].Color
		}
	}

	finalFrame := make(Frame, 0, len(frameAsMap))

	for _, pixel := range frameAsMap {
		finalFrame = append(finalFrame, *pixel)
	}

	return finalFrame
}
