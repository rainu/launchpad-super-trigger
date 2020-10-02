package pad

import (
	"fmt"
	"github.com/rainu/launchpad"
	"strings"
)

var (
	ColorOff = LaunchpadDependColor{
		S:   LaunchpadSColor{0, 0},
		MK2: LaunchpadMK2RGBColor{0, 0, 0},
	}

	ColorGreen = LaunchpadDependColor{
		S:   LaunchpadSColor{3, 0},
		MK2: LaunchpadMK2RGBColor{0, 63, 0},
	}
	ColorOrange = LaunchpadDependColor{
		S:   LaunchpadSColor{2, 3},
		MK2: LaunchpadMK2RGBColor{63, 31, 0},
	}
	ColorRed = LaunchpadDependColor{
		S:   LaunchpadSColor{0, 3},
		MK2: LaunchpadMK2RGBColor{63, 0, 0},
	}
)

type Color interface {
	Light(pad Lighter, x, y int) error
	Text(pad Lighter) launchpad.ScrollingTextBuilder
	TextLoop(pad Lighter) launchpad.ScrollingTextBuilder
	Ordinal() string
}

//LaunchpadSColor represents a color for the "Launchpad S"
type LaunchpadSColor struct {
	Green int
	Red   int
}

func (c LaunchpadSColor) Light(pad Lighter, x, y int) error {
	return pad.Light(x, y, launchpad.ColorS{Green: c.Green, Red: c.Red})
}

func (c LaunchpadSColor) Text(pad Lighter) launchpad.ScrollingTextBuilder {
	return pad.Text(launchpad.ColorS{Green: c.Green, Red: c.Red})
}

func (c LaunchpadSColor) TextLoop(pad Lighter) launchpad.ScrollingTextBuilder {
	return pad.TextLoop(launchpad.ColorS{Green: c.Green, Red: c.Red})
}

func (c LaunchpadSColor) Ordinal() string {
	return fmt.Sprintf("%x%x", c.Red, c.Green)
}

//LaunchpadMK2RGBColor represents a rgb color for the "Launchpad MK2"
type LaunchpadMK2RGBColor struct {
	Red   int
	Green int
	Blue  int
}

func (c LaunchpadMK2RGBColor) Light(pad Lighter, x, y int) error {
	return pad.Light(x, y, launchpad.ColorMK2RGB{Red: c.Red, Green: c.Green, Blue: c.Blue})
}

func (c LaunchpadMK2RGBColor) Text(pad Lighter) launchpad.ScrollingTextBuilder {
	return pad.Text(launchpad.ColorMK2RGB{Red: c.Red, Green: c.Green, Blue: c.Blue})
}

func (c LaunchpadMK2RGBColor) TextLoop(pad Lighter) launchpad.ScrollingTextBuilder {
	return pad.TextLoop(launchpad.ColorMK2RGB{Red: c.Red, Green: c.Green, Blue: c.Blue})
}

func (c LaunchpadMK2RGBColor) Ordinal() string {
	return fmt.Sprintf("%x%x%x", c.Red, c.Green, c.Blue)
}

//LaunchpadMK2CodeColor represents a color (code) for the "Launchpad MK2"
type LaunchpadMK2CodeColor struct {
	Code int
}

func (c LaunchpadMK2CodeColor) Light(pad Lighter, x, y int) error {
	return pad.Light(x, y, launchpad.ColorMK2{Code: c.Code})
}

func (c LaunchpadMK2CodeColor) Text(pad Lighter) launchpad.ScrollingTextBuilder {
	return pad.Text(launchpad.ColorMK2{Code: c.Code})
}

func (c LaunchpadMK2CodeColor) TextLoop(pad Lighter) launchpad.ScrollingTextBuilder {
	return pad.TextLoop(launchpad.ColorMK2{Code: c.Code})
}

func (c LaunchpadMK2CodeColor) Ordinal() string {
	return fmt.Sprintf("%x", c.Code)
}

//LaunchpadDependColor contains multiple colors, one for each supported launchpad type (S, MK2)
type LaunchpadDependColor struct {
	S   LaunchpadSColor
	MK2 LaunchpadMK2RGBColor
}

func (c LaunchpadDependColor) Light(pad Lighter, x, y int) error {
	if isLaunchpadS(pad) {
		return c.S.Light(pad, x, y)
	}

	return c.MK2.Light(pad, x, y)
}

func (c LaunchpadDependColor) Text(pad Lighter) launchpad.ScrollingTextBuilder {
	if isLaunchpadS(pad) {
		return c.S.Text(pad)
	}

	return c.MK2.Text(pad)
}

func (c LaunchpadDependColor) TextLoop(pad Lighter) launchpad.ScrollingTextBuilder {
	if isLaunchpadS(pad) {
		return c.S.TextLoop(pad)
	}

	return c.MK2.TextLoop(pad)
}

func (c LaunchpadDependColor) Ordinal() string {
	return c.S.Ordinal() + c.MK2.Ordinal()
}

func isLaunchpadS(pad Lighter) bool {
	return strings.Contains(pad.Name(), "Launchpad S")
}
