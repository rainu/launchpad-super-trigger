package mqtt

import (
	MQTT "github.com/goiiot/libmqtt"
	"go.uber.org/zap"
)

type logConnectionListener struct {
}

func (l *logConnectionListener) OnConnect(client MQTT.Client) {
	zap.L().Info("MQTT connection established broker.")
}

func (l *logConnectionListener) OnConnectionLost(client MQTT.Client, err error) {
	zap.L().Warn("Connection lost to broker.", zap.Error(err))
}
