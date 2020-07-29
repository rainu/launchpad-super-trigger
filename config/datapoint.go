package config

import "strings"

type Datapoint string

func (d Datapoint) IsValid() bool {
	split := strings.Split(string(d), ".")
	return len(split) == 2
}

func (d Datapoint) Path() string {
	return string(d)
}

func (d Datapoint) Sensor() string {
	return strings.Split(string(d), ".")[0]
}

func (d Datapoint) Name() string {
	return strings.Split(string(d), ".")[1]
}
