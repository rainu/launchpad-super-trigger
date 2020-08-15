package sensor

import (
	"fmt"
	"github.com/boltdb/bolt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/sensor"
	"github.com/rainu/launchpad-super-trigger/sensor/data_extractor"
	"github.com/rainu/launchpad-super-trigger/sensor/store"
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

func BuildSensors(generalSettings config.General, sensors config.Sensors, mqttConnections map[string]MQTT.Client) map[string]Sensor {
	result := map[string]Sensor{}

	buildMqttSensors(result, generalSettings, sensors.Mqtt, mqttConnections)
	buildRestSensors(result, generalSettings, sensors.Rest)
	buildCommandSensors(result, generalSettings, sensors.Command)

	return result
}

var usedBoltDB *bolt.DB

func generateStore(generalSettings config.General, sensorName string) store.Store {
	var result store.Store

	if generalSettings.SensorStore == "" {
		result = &store.MemoryStore{}
	} else {
		if usedBoltDB == nil {
			var err error
			usedBoltDB, err = bolt.Open(generalSettings.SensorStore, 0600, nil)
			if err != nil {
				panic(fmt.Errorf("could not open boltDB: %w", err))
			}
		}

		result = &store.BoltDBStore{
			BoltDB: usedBoltDB,
			Bucket: "sensor_data",
			Key:    sensorName,
		}
	}

	if generalSettings.CompressSensorData {
		result = &store.GzipCompressed{
			Delegate: result,
		}
	}

	return result
}
