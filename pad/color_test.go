package pad

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColor_Ordinal(t *testing.T) {
	colorOrdinals := map[byte]bool{}

	colorOrdinals[ColorOff.Ordinal()] = true
	colorOrdinals[ColorDimRed.Ordinal()] = true
	colorOrdinals[ColorNormalRed.Ordinal()] = true
	colorOrdinals[ColorHighRed.Ordinal()] = true
	colorOrdinals[ColorDimGreen.Ordinal()] = true
	colorOrdinals[ColorNormalGreen.Ordinal()] = true
	colorOrdinals[ColorHighGreen.Ordinal()] = true
	colorOrdinals[ColorDimYellow.Ordinal()] = true
	colorOrdinals[ColorNormalYellow.Ordinal()] = true
	colorOrdinals[ColorHighYellow.Ordinal()] = true
	colorOrdinals[ColorDimOrange.Ordinal()] = true
	colorOrdinals[ColorNormalOrange.Ordinal()] = true
	colorOrdinals[ColorHighOrange.Ordinal()] = true
	colorOrdinals[ColorDimLightGreen.Ordinal()] = true
	colorOrdinals[ColorLightGreen.Ordinal()] = true
	colorOrdinals[ColorHighLightGreen.Ordinal()] = true

	assert.Equal(t, 16, len(colorOrdinals))
}
