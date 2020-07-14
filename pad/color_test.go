package pad

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColor_Ordinal(t *testing.T) {
	colorOrdinals := map[byte]bool{}

	colorOrdinals[Off.Ordinal()] = true
	colorOrdinals[DimRed.Ordinal()] = true
	colorOrdinals[NormalRed.Ordinal()] = true
	colorOrdinals[HighRed.Ordinal()] = true
	colorOrdinals[DimGreen.Ordinal()] = true
	colorOrdinals[NormalGreen.Ordinal()] = true
	colorOrdinals[HighGreen.Ordinal()] = true
	colorOrdinals[DimYellow.Ordinal()] = true
	colorOrdinals[NormalYellow.Ordinal()] = true
	colorOrdinals[HighYellow.Ordinal()] = true
	colorOrdinals[DimOrange.Ordinal()] = true
	colorOrdinals[NormalOrange.Ordinal()] = true
	colorOrdinals[HighOrange.Ordinal()] = true
	colorOrdinals[DimLightGreen.Ordinal()] = true
	colorOrdinals[LightGreen.Ordinal()] = true
	colorOrdinals[HighLightGreen.Ordinal()] = true

	assert.Equal(t, 16, len(colorOrdinals))
}
