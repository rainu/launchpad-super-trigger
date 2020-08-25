package mqtt

import (
	"context"
	MQTT "github.com/goiiot/libmqtt"
	"github.com/rainu/launchpad-super-trigger/config"
	"go.uber.org/zap"
)

type Client interface {
	Connect(ctx context.Context) error
	Publish(ctx context.Context, topic string, qos byte, retained bool, payload []byte) error
	Subscribe(topic string, qos byte, clb func(string, byte, []byte))
	AddListener(cl ConnectionListener)
}

func BuildMqttConnection(parsedConfig *config.Config) map[string]Client {
	result := map[string]Client{}

	for name, connection := range parsedConfig.Connections.MQTT {
		result[name] = buildMqttConnection(connection)
	}

	return result
}

func buildMqttConnection(connection config.MQTTConnection) Client {
	result := &mqttClient{
		conObserver: connectionObserver{},
	}

	client, err := MQTT.NewClient(
		MQTT.WithClientID(connection.ClientId),
		MQTT.WithIdentity(connection.Username, connection.Password),
		MQTT.WithCleanSession(true),
		MQTT.WithAutoReconnect(true),
	)
	if err != nil {
		zap.L().Fatal("Error build connecting to mqtt broker: %s", zap.Error(err))
	}

	result.broker = connection.Broker
	result.client = client
	result.conObserver.AddListener(&logConnectionListener{}) //for logging purposes

	return result
}
