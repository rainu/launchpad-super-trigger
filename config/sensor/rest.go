package sensor

import (
	"encoding/base64"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/sensor"
	"io"
	"net/http"
	"os"
	"strings"
)

func buildRestSensors(sensors map[string]Sensor, generalSettings config.General, restSensors map[string]config.RestSensor) {
	for sensorName, restSensor := range restSensors {
		s := &sensor.Rest{
			HttpClient:   http.DefaultClient,
			Method:       restSensor.Method,
			Url:          restSensor.URL,
			Header:       restSensor.Header,
			Interval:     restSensor.Interval,
			MessageStore: generateStore(generalSettings, sensorName),
		}

		if restSensor.BodyRaw != "" {
			s.Body = rawBody(restSensor.BodyRaw)
		} else if restSensor.BodyB64 != "" {
			s.Body = b64Body(restSensor.BodyB64)
		} else if restSensor.BodyPath != "" {
			s.Body = fileBody(restSensor.BodyPath)
		}

		sensors[sensorName] = Sensor{
			Sensor:     s,
			Extractors: buildExtractors(restSensor.DataPoints),
		}
	}
}

func rawBody(body string) func() io.Reader {
	return func() io.Reader {
		return strings.NewReader(body)
	}
}

func b64Body(b64Body string) func() io.Reader {
	return func() io.Reader {
		return base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64Body))
	}
}

func fileBody(bodyPath string) func() io.Reader {
	return func() io.Reader {
		file, err := os.Open(bodyPath)
		if err != nil {
			panic(err)
		}

		return file
	}
}
