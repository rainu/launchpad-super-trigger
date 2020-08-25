package mqtt

import (
	"context"
	MQTT "github.com/goiiot/libmqtt"
)

type mqttClient struct {
	broker      string
	client      MQTT.Client
	conObserver connectionObserver
}

type connectionListener struct {
	errChan chan error
}

func (c *connectionListener) OnConnect(client MQTT.Client) {
	c.errChan <- nil
}

func (c *connectionListener) OnConnectionLost(client MQTT.Client, err error) {
	c.errChan <- err
}

func (m *mqttClient) Connect(ctx context.Context) error {
	errChan := make(chan error, 1)
	defer close(errChan)

	listener := &connectionListener{errChan}
	defer func() {
		//we can remove our listener -> the only purpose was to wait for connection!
		m.conObserver.RemoveListener(listener)
	}()

	m.conObserver.AddListener(listener)
	err := m.client.ConnectServer(m.broker, MQTT.WithConnHandleFunc(m.conObserver.OnConnectHandler()))
	if err != nil {
		return err
	}

	//wait until connection established
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *mqttClient) Subscribe(topic string, qos byte, clb func(string, byte, []byte)) {
	m.client.HandleTopic(topic, func(client MQTT.Client, topic string, qos MQTT.QosLevel, msg []byte) {
		clb(topic, qos, msg)
	})
	m.client.Subscribe(&MQTT.Topic{
		Name: topic,
		Qos:  qos,
	})
}

func (m *mqttClient) Publish(ctx context.Context, topic string, qos byte, retained bool, payload []byte) error {
	m.client.Publish(&MQTT.PublishPacket{
		TopicName: topic,
		Qos:       qos,
		IsRetain:  retained,
		Payload:   payload,
	})
	return nil
}
