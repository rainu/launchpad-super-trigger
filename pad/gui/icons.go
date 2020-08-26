package gui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
)

type colorIcon struct {
	name    string
	color   color.Color
	content []byte
}

func (r colorIcon) Name() string {
	return r.name
}

func (r colorIcon) Content() []byte {
	if r.content == nil {
		img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{1, 1}})
		img.Set(0, 0, r.color)

		buf := bytes.Buffer{}
		png.Encode(&buf, img)

		r.content = buf.Bytes()
	}

	return r.content
}

type launchpadColor struct {
	r int
	g int
}

func (l launchpadColor) RGBA() (uint32, uint32, uint32, uint32) {
	return color.RGBA{uint8(l.r) * 64, uint8(l.g) * 64, 0, 0xff}.RGBA()
}

var loadedIcons = map[string]*colorIcon{}

func getColorIcon(r, g int) *colorIcon {
	key := fmt.Sprintf("%d,%d", r, g)
	if loadedIcons[key] == nil {
		loadedIcons[key] = &colorIcon{
			name:  key,
			color: launchpadColor{r: r, g: g},
		}
	}

	return loadedIcons[key]
}
