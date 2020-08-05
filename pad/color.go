package pad

import "github.com/rainu/launchpad"

var (
	ColorOff = Color{0, 0}

	ColorDimRed    = Color{0, 1}
	ColorNormalRed = Color{0, 2}
	ColorHighRed   = Color{0, 3}

	ColorDimGreen    = Color{1, 0}
	ColorNormalGreen = Color{2, 0}
	ColorHighGreen   = Color{3, 0}

	ColorDimYellow    = Color{1, 1}
	ColorNormalYellow = Color{2, 2}
	ColorHighYellow   = Color{3, 3}

	ColorDimOrange    = Color{1, 2}
	ColorNormalOrange = Color{1, 3}
	ColorHighOrange   = Color{2, 3}

	ColorDimLightGreen  = Color{2, 1}
	ColorLightGreen     = Color{3, 1}
	ColorHighLightGreen = Color{3, 2}

	AllColors = []Color{
		ColorDimRed,
		ColorNormalRed,
		ColorHighRed,
		ColorDimGreen,
		ColorNormalGreen,
		ColorHighGreen,
		ColorDimYellow,
		ColorNormalYellow,
		ColorHighYellow,
		ColorDimOrange,
		ColorNormalOrange,
		ColorHighOrange,
		ColorDimLightGreen,
		ColorLightGreen,
		ColorHighLightGreen,
	}
)

type Color struct {
	Green int
	Red   int
}

func (c Color) Light(pad Lighter, x, y int) error {
	return pad.Light(x, y, c.Green, c.Red)
}

func (c Color) Text(pad Lighter) launchpad.ScrollingTextBuilder {
	return pad.Text(c.Green, c.Red)
}

func (c Color) TextLoop(pad Lighter) launchpad.ScrollingTextBuilder {
	return pad.TextLoop(c.Green, c.Red)
}

func (c Color) Ordinal() byte {
	return byte(4*c.Red + c.Green)
}
