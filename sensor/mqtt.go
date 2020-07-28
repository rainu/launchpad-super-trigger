package sensor

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"sync"
)

type MQTT struct {
	callbackHandler

	Client mqtt.Client
	Topic  string
	QOS    byte

	running     bool
	mux         sync.RWMutex
	lastMessage []byte
}

func (m *MQTT) Run(ctx context.Context) error {
	if m.running {
		return fmt.Errorf("listerner is already running")
	}

	m.running = true
	m.Client.Subscribe(m.Topic, m.QOS, m.handleMessage)

	//wait until context closed
	<-ctx.Done()

	return nil
}

func (m *MQTT) LastMessage() []byte {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.lastMessage
}

func (m *MQTT) handleMessage(client mqtt.Client, message mqtt.Message) {
	zap.L().Debug(fmt.Sprintf("Mqtt message received: %s", message.Topic()))
	message.Ack()

	m.mux.Lock()
	m.lastMessage = message.Payload()
	m.mux.Unlock()

	m.callbackHandler.Call(m)
}
