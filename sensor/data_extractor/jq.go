package data_extractor

import (
	"encoding/json"
	"fmt"
	"github.com/itchyny/gojq"
	"strings"
)

type Jq struct {
	Query    string
	compiled *gojq.Code
}

func (j *Jq) Extract(data []byte) ([]byte, error) {
	if j.compiled == nil {
		query, err := gojq.Parse(j.Query)
		if err != nil {
			return nil, err
		}

		j.compiled, err = gojq.Compile(query)
		if err != nil {
			return nil, err
		}
	}

	sb := strings.Builder{}

	parsed := map[string]interface{}{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	iter := j.compiled.Run(parsed)

	isFirst := true
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}

		if isFirst {
			isFirst = false
			sb.WriteString(fmt.Sprintf("%#v", v))
		} else {
			sb.WriteString(fmt.Sprintf("\n%#v", v))
		}
	}

	return []byte(sb.String()), nil
}
