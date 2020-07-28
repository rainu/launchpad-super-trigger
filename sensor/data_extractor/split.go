package data_extractor

import (
	"errors"
	"strings"
)

type Split struct {
	Separator string
	Index     int
}

func (s Split) Extract(data []byte) ([]byte, error) {
	split := strings.Split(string(data), s.Separator)
	if s.Index > len(split) {
		return nil, errors.New("index out of bounds")
	}

	return []byte(split[s.Index]), nil
}
