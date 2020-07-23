package config

import (
	"errors"
	"regexp"
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

type Coordinates string

var coordRangeRegex = regexp.MustCompile(`[0-7]-[0-7]`)

func (c Coordinates) Coordinates() ([][]int, error) {
	x, y, err := Coordinate(c).Coordinate()
	if err == nil {
		//its a single coordinate!
		return [][]int{{x, y}}, nil
	}

	split := strings.Split(string(c), ",")
	if len(split) != 2 {
		return nil, errors.New("syntax error")
	}

	rawX := split[0]
	rawY := split[1]

	xStart, xUntil, err := convertRange(rawX)
	if err != nil {
		return nil, err
	}
	yStart, yUntil, err := convertRange(rawY)
	if err != nil {
		return nil, err
	}

	result := make([][]int, 0, 8*8)

	for x := xStart; x <= xUntil; x++ {
		for y := yStart; y <= yUntil; y++ {
			result = append(result, []int{x, y})
		}
	}

	return result, nil
}

func convertRange(raw string) (int, int, error) {
	start := 0
	until := 0

	i, err := strconv.Atoi(raw)
	if err == nil {
		if i < 0 || i > 7 {
			return -1, -1, errors.New("invalid value")
		}
		start = i
		until = i
	} else if coordRangeRegex.MatchString(raw) {
		split := strings.Split(raw, "-")
		start, _ = strconv.Atoi(split[0])
		until, _ = strconv.Atoi(split[1])
	} else {
		return -1, -1, errors.New("syntax error")
	}

	return start, until, nil
}
