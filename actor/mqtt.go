package actor

import (
	"bytes"
	"context"
	"io"
)

type MQTTPublisher interface {
	Publish(ctx context.Context, topic string, qos byte, retained bool, payload []byte) error
}

type Mqtt struct {
	Client   MQTTPublisher
	Topic    string
	QOS      byte
	Retained bool
	Body     func() io.Reader
}

func (m *Mqtt) Do(ctx Context) error {
	byteBuff := bytes.NewBuffer(make([]byte, 0, 8192))
	_, err := byteBuff.ReadFrom(m.Body())
	if err != nil {
		return err
	}

	return m.Client.Publish(ctx.Context, m.Topic, m.QOS, m.Retained, byteBuff.Bytes())
}
