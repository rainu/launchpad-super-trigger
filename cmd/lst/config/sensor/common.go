package sensor

import (
	"github.com/rainu/launchpad-super-trigger/sensor"
	"github.com/rainu/launchpad-super-trigger/sensor/data_extractor"
	"sync"
)

type Sensor struct {
	Sensor     sensor.Sensor
	Extractors map[string]data_extractor.Extractor
}

type DataListener interface {
	OnData(dataName string, data []byte, err error)
}

type DataObserver struct {
	Extractors    map[string]data_extractor.Extractor
	listener      map[interface{}]DataListener
	listenerMutex sync.RWMutex
}

func (d *DataObserver) AddListener(dl DataListener) {
	d.listenerMutex.Lock()
	defer d.listenerMutex.Unlock()

	d.listener[dl] = dl
}

func (d *DataObserver) RemoveListener(dl DataListener) {
	d.listenerMutex.Lock()
	defer d.listenerMutex.Unlock()

	delete(d.listener, dl)
}
