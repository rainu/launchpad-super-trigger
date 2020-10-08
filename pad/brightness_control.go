package pad

import (
	"github.com/rainu/launchpad"
	"go.uber.org/zap"
	"sync"
)

type BrightnessLevel byte

//this middleware is responsible for controlling the brightness
type brightnessMiddleware struct {
	delegate   Launchpad
	brightness BrightnessLevel
	colorGrid  map[int]map[int]launchpad.Color
	mux        sync.Mutex
}

func (r *realLaunchpad) SetBrightness(b BrightnessLevel) {
	//this is empty because the brightnessMiddleware contains the implementation
}

func (r *realLaunchpad) GetBrightness() BrightnessLevel {
	//this is empty because the brightnessMiddleware contains the implementation
	return 0
}

func (b *brightnessMiddleware) GetBrightness() BrightnessLevel {
	return b.brightness
}

func (b *brightnessMiddleware) SetBrightness(b2 BrightnessLevel) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.brightness = b2

	//apply changes
	if b.colorGrid != nil {
		for x, yGrid := range b.colorGrid {
			for y, color := range yGrid {
				b.delegate.Light(x, y, b.adjustColor(color))
			}
		}
	}
}

func (b *brightnessMiddleware) Light(x, y int, color launchpad.Color) error {
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.colorGrid == nil {
		b.colorGrid = map[int]map[int]launchpad.Color{}
	}
	if b.colorGrid[x] == nil {
		b.colorGrid[x] = map[int]launchpad.Color{}
	}
	b.colorGrid[x][y] = color

	return b.delegate.Light(x, y, b.adjustColor(color))
}

func (b *brightnessMiddleware) Text(color launchpad.Color) launchpad.ScrollingTextBuilder {
	b.mux.Lock()
	defer b.mux.Unlock()

	return b.delegate.Text(b.adjustColor(color))
}

func (b *brightnessMiddleware) TextLoop(color launchpad.Color) launchpad.ScrollingTextBuilder {
	b.mux.Lock()
	defer b.mux.Unlock()

	return b.delegate.TextLoop(b.adjustColor(color))
}

func (b *brightnessMiddleware) adjustColor(color launchpad.Color) launchpad.Color {
	// only if brightness should change
	if b.brightness != 100 {
		rgbColor, ok := color.(launchpad.ColorMK2RGB)
		if !ok {
			zap.L().Debug("Brightness control is not supported for that color schema.")
		} else {
			return launchpad.ColorMK2RGB{
				Red:   rgbColor.Red * int(b.brightness) / 100,
				Green: rgbColor.Green * int(b.brightness) / 100,
				Blue:  rgbColor.Blue * int(b.brightness) / 100,
			}
		}
	}

	return color
}

func (b *brightnessMiddleware) Name() string {
	return b.delegate.Name()
}

func (b *brightnessMiddleware) Clear() error {
	b.mux.Lock()
	defer b.mux.Unlock()

	//forget all colors
	b.colorGrid = map[int]map[int]launchpad.Color{}

	return b.delegate.Clear()
}

func (b *brightnessMiddleware) Close() error {
	return b.delegate.Close()
}

func (b *brightnessMiddleware) ListenToHits() (<-chan launchpad.Hit, error) {
	return b.delegate.ListenToHits()
}

func (b *brightnessMiddleware) IsHealthy() bool {
	return b.delegate.IsHealthy()
}
