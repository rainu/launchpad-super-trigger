package data_extractor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGJSON_Extract_Happy(t *testing.T) {
	toTest := GJSON{"name.first"}

	json := `{"name":{"first":"Max","last":"Mustermann"},"age":13}`

	extract, err := toTest.Extract([]byte(json))

	assert.NoError(t, err)
	assert.Equal(t, "Max", string(extract))
}

func TestGJSON_Extract_Happy2(t *testing.T) {
	toTest := GJSON{"age"}

	json := `{"name":{"first":"Max","last":"Mustermann"},"age":13}`

	extract, err := toTest.Extract([]byte(json))

	assert.NoError(t, err)
	assert.Equal(t, "13", string(extract))
}

func TestGJSON_Extract_Fail(t *testing.T) {
	toTest := GJSON{"DoesNotExists"}

	json := `{"name":{"first":"Max","last":"Mustermann"},"age":13}`

	extract, err := toTest.Extract([]byte(json))

	assert.NoError(t, err)
	assert.Equal(t, "", string(extract))
}
