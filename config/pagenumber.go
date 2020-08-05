package config

import (
	"regexp"
	"strconv"
	"strings"
)

type PageNumber string

var pageNumberRegex = regexp.MustCompile("^[01]{8}$")

func (p PageNumber) isValid() bool {
	n := p.AsInt()
	return n >= 0 && n <= 255
}

func (p PageNumber) AsInt() int {
	if pageNumberRegex.Match([]byte(p)) {
		n, err := strconv.ParseInt(p.reverse(), 2, 32)
		if err != nil {
			return -1
		}
		return int(n)
	}

	n, err := strconv.ParseInt(string(p), 10, 32)
	if err != nil {
		return -1
	}
	return int(n)
}

func (p PageNumber) reverse() string {
	sb := strings.Builder{}

	for i := len(p) - 1; i >= 0; i-- {
		sb.Write([]byte{[]byte(p)[i]})
	}

	return sb.String()
}
