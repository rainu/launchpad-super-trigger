package sensor

import (
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/sensor/data_extractor"
)

func buildExtractors(dp config.DataPoints) map[string]data_extractor.Extractor {
	result := map[string]data_extractor.Extractor{}

	buildGjsonExtractors(result, dp.Gjson)
	buildGojqExtractors(result, dp.Gojq)
	buildSplitExtractors(result, dp.Split)

	return result
}

func buildGjsonExtractors(extractors map[string]data_extractor.Extractor, gjson map[string]string) {
	for name, gjsonPath := range gjson {
		extractors[name] = data_extractor.GJSON{
			Path: gjsonPath,
		}
	}
}

func buildGojqExtractors(extractors map[string]data_extractor.Extractor, gojq map[string]string) {
	for name, query := range gojq {
		extractors[name] = &data_extractor.Jq{
			Query: query,
		}
	}
}

func buildSplitExtractors(extractors map[string]data_extractor.Extractor, split map[string]config.SplitDataPoint) {
	for name, splitDP := range split {
		extractors[name] = data_extractor.Split{
			Separator: splitDP.Separator,
			Index:     splitDP.Index,
		}
	}
}
