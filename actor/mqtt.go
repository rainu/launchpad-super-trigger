package actor

import (
	"bytes"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io"
)

type MqttActionHandler struct {
	Client   mqtt.Client
	Topic    string
	QOS      byte
	Retained bool
	Body     func() io.Reader
}

func (m *MqttActionHandler) Do(ctx Context) error {
	byteBuff := bytes.NewBuffer(make([]byte, 0, 8192))
	_, err := byteBuff.ReadFrom(m.Body())
	if err != nil {
		return err
	}

	token := m.Client.Publish(m.Topic, m.QOS, m.Retained, byteBuff.Bytes())

	waitChan := make(chan error, 1)

	go func() {
		token.Wait()
		waitChan <- token.Error()
	}()

	select {
	case err := <-waitChan:
		return err
	case <-ctx.Context.Done():
		return ctx.Context.Err()
	}
}
