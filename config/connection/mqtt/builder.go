package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/config"
)

type Client interface {
	MQTT.Client
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

	opts := MQTT.NewClientOptions()
	opts.AddBroker(connection.Broker)

	opts.SetAutoReconnect(true)
	opts.SetClientID(connection.ClientId)
	opts.SetUsername(connection.Username)
	opts.SetPassword(connection.Password)
	opts.SetOnConnectHandler(result.conObserver.OnConnectHandler())
	opts.SetConnectionLostHandler(result.conObserver.ConnectionLostHandler())

	result.client = MQTT.NewClient(opts)
	result.conObserver.AddListener(&logConnectionListener{}) //for logging purposes

	return result
}
