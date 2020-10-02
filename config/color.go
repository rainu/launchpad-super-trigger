package config

import (
	"errors"
	"github.com/rainu/launchpad-super-trigger/pad"
	"strconv"
	"strings"
)

type Color string

func (c Color) Color() (pad.Color, error) {
	split := strings.Split(string(c), ",")

	switch len(split) {
	case 2:
		return c.handleLaunchpadSColor(split)
	case 3:
		return c.handleLaunchpadMK2Color(split)
	}

	return nil, errors.New("syntax error")
}

func (c Color) handleLaunchpadSColor(split []string) (pad.Color, error) {
	r, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, err
	}
	g, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, err
	}

	if r < 0 || r > 3 {
		return nil, errors.New("invalid red value")
	}
	if g < 0 || g > 3 {
		return nil, errors.New("invalid green value")
	}

	return pad.LaunchpadSColor{
		Green: g,
		Red:   r,
	}, nil
}

func (c Color) handleLaunchpadMK2Color(split []string) (pad.Color, error) {
	r, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, err
	}
	g, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, err
	}
	b, err := strconv.Atoi(split[2])
	if err != nil {
		return nil, err
	}

	if r < 0 || r > 63 {
		return nil, errors.New("invalid red value")
	}
	if g < 0 || g > 63 {
		return nil, errors.New("invalid green value")
	}
	if b < 0 || b > 63 {
		return nil, errors.New("invalid blue value")
	}

	return pad.LaunchpadMK2RGBColor{
		Green: g,
		Red:   r,
		Blue:  b,
	}, nil
}
