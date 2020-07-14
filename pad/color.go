package pad

var (
	Off = Color{0, 0}

	DimRed    = Color{0, 1}
	NormalRed = Color{0, 2}
	HighRed   = Color{0, 3}

	DimGreen    = Color{1, 0}
	NormalGreen = Color{2, 0}
	HighGreen   = Color{3, 0}

	DimYellow    = Color{1, 1}
	NormalYellow = Color{2, 2}
	HighYellow   = Color{3, 3}

	DimOrange    = Color{1, 2}
	NormalOrange = Color{1, 3}
	HighOrange   = Color{2, 3}

	DimLightGreen  = Color{2, 1}
	LightGreen     = Color{3, 1}
	HighLightGreen = Color{3, 2}
)

type Color struct {
	Green int
	Red   int
}

func (c Color) Light(pad Lighter, x, y int) error {
	return pad.Light(x, y, c.Green, c.Red)
}

func (c Color) Ordinal() byte {
	return byte(4*c.Red + c.Green)
}
