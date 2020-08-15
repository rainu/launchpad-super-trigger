package sensor

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/sensor/store"
	"go.uber.org/zap"
)

type MQTT struct {
	callbackHandler

	Client       mqtt.Client
	Topic        string
	QOS          byte
	MessageStore store.Store

	running bool
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
	data, err := m.MessageStore.Get()
	if err != nil {
		zap.L().Error("Could not get message from message store!", zap.Error(err))
		return nil
	}

	return data
}

func (m *MQTT) handleMessage(client mqtt.Client, message mqtt.Message) {
	zap.L().Debug(fmt.Sprintf("Mqtt message received: %s", message.Topic()))
	message.Ack()

	if err := m.MessageStore.Set(message.Payload()); err != nil {
		zap.L().Error("Could not save message into message store!", zap.Error(err))
	}
	m.callbackHandler.Call(m)
}
