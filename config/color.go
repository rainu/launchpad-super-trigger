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
	if len(split) != 2 {
		return pad.Color{}, errors.New("syntax error")
	}

	r, err := strconv.Atoi(split[0])
	if err != nil {
		return pad.Color{}, err
	}
	g, err := strconv.Atoi(split[1])
	if err != nil {
		return pad.Color{}, err
	}

	if r < 0 || r > 3 {
		return pad.Color{}, errors.New("invalid red value")
	}
	if g < 0 || g > 3 {
		return pad.Color{}, errors.New("invalid green value")
	}

	return pad.Color{
		Green: g,
		Red:   r,
	}, nil
}
