package data_extractor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJq_Extract_Happy_String(t *testing.T) {
	toTest := Jq{Query: `.name.first`}

	json := `{"name":{"first":"Max","last":"Mustermann"},"age":13}`

	extract, err := toTest.Extract([]byte(json))

	assert.NoError(t, err)
	assert.Equal(t, `"Max"`, string(extract))
}

func TestJq_Extract_Happy_Int(t *testing.T) {
	toTest := Jq{Query: `.age`}

	json := `{"name":{"first":"Max","last":"Mustermann"},"age":13}`

	extract, err := toTest.Extract([]byte(json))

	assert.NoError(t, err)
	assert.Equal(t, `13`, string(extract))
}

func TestJq_Extract_Happy_Arithmetic(t *testing.T) {
	toTest := Jq{Query: `.cur * .max / 100`}

	json := `{"cur": 5, "max": 500}`

	extract, err := toTest.Extract([]byte(json))

	assert.NoError(t, err)
	assert.Equal(t, `25`, string(extract))
}

func TestJq_Extract_Fail(t *testing.T) {
	toTest := Jq{Query: `QuerySyntaxError`}

	json := `{"name":{"first":"Max","last":"Mustermann"},"age":13}`

	_, err := toTest.Extract([]byte(json))

	assert.Error(t, err)
}
