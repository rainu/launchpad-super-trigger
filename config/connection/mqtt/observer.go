package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
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
	listener      []ConnectionListener
}

func (c *connectionObserver) AddListener(cl ConnectionListener) {
	c.listenerMutex.Lock()
	defer c.listenerMutex.Unlock()

	c.listener = append(c.listener, cl)
}

func (c *connectionObserver) OnConnectHandler() MQTT.OnConnectHandler {
	return func(client MQTT.Client) {
		c.listenerMutex.RLock()
		c.eventMutex.Lock()
		defer c.listenerMutex.RUnlock()
		defer c.eventMutex.Unlock()

		for _, listener := range c.listener {
			listener.OnConnect(client)
		}
	}
}

func (c *connectionObserver) ConnectionLostHandler() MQTT.ConnectionLostHandler {
	return func(client MQTT.Client, err error) {
		c.listenerMutex.RLock()
		c.eventMutex.Lock()
		defer c.listenerMutex.RUnlock()
		defer c.eventMutex.Unlock()

		for _, listener := range c.listener {
			listener.OnConnectionLost(client, err)
		}
	}
}
