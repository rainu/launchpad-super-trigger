package gfx

import "github.com/rainu/launchpad-super-trigger/pad"

type FramePixel struct {
	X     int
	Y     int
	Color pad.Color
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

func joinPixels(pixels ...FramePixel) []FramePixel {
	alreadyKnownPositions := map[coord]bool{}
	joined := make([]FramePixel, 0, len(pixels))

	for _, pixel := range pixels {
		if !alreadyKnownPositions[coord{pixel.X, pixel.Y}] {
			alreadyKnownPositions[coord{pixel.X, pixel.Y}] = true
			joined = append(joined, pixel)
		}
	}

	return joined
}
