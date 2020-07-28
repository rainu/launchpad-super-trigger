package data_extractor

import (
	"fmt"
	"github.com/tidwall/gjson"
)

type GJSON struct {
	Path string
}

func (g GJSON) Extract(data []byte) ([]byte, error) {
	result := gjson.GetBytes(data, g.Path)
	if result.Exists() {
		s := fmt.Sprintf("%v", result.Value())
		return []byte(s), nil
	}

	return []byte{}, nil
}
