package config

import (
	"errors"
	"strconv"
	"strings"
)

type Coordinate string

func (c Coordinate) Coordinate() (int, int, error) {
	split := strings.Split(string(c), ",")
	if len(split) != 2 {
		return -1, -1, errors.New("syntax error")
	}

	x, err := strconv.Atoi(split[0])
	if err != nil {
		return -1, -1, err
	}
	y, err := strconv.Atoi(split[1])
	if err != nil {
		return -1, -1, err
	}

	if x < 0 || x > 7 {
		return -1, -1, errors.New("invalid x value")
	}
	if y < 0 || y > 7 {
		return -1, -1, errors.New("invalid y value")
	}

	return x, y, nil
}
