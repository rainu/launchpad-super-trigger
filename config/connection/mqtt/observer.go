package mqtt

import (
	MQTT "github.com/goiiot/libmqtt"
	"sync"
)

func (m *mqttClient) AddListener(cl ConnectionListener) {
	m.conObserver.AddListener(cl)
}

type ConnectionListener interface {
	OnConnect(MQTT.Client)
	OnConnectionLost(MQTT.Client, error)
}

type connectionObserver struct {
	listenerMutex sync.RWMutex
	eventMutex    sync.Mutex
	listener      map[interface{}]ConnectionListener
}

func (c *connectionObserver) AddListener(cl ConnectionListener) {
	c.listenerMutex.Lock()
	defer c.listenerMutex.Unlock()

	if c.listener == nil {
		c.listener = map[interface{}]ConnectionListener{}
	}

	c.listener[cl] = cl
}
func (c *connectionObserver) RemoveListener(cl ConnectionListener) {
	c.listenerMutex.Lock()
	defer c.listenerMutex.Unlock()

	delete(c.listener, cl)
}

func (c *connectionObserver) OnConnectHandler() MQTT.ConnHandleFunc {
	return func(client MQTT.Client, server string, code byte, err error) {
		c.listenerMutex.RLock()
		c.eventMutex.Lock()
		defer c.listenerMutex.RUnlock()
		defer c.eventMutex.Unlock()

		if code == MQTT.CodeSuccess {
			for _, listener := range c.listener {
				listener.OnConnect(client)
			}
		} else {
			for _, listener := range c.listener {
				listener.OnConnectionLost(client, err)
			}
		}
	}
}
