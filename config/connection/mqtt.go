package connection

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rainu/launchpad-super-trigger/config"
	"go.uber.org/zap"
)

func BuildMqttConnection(parsedConfig *config.Config) map[string]MQTT.Client {
	result := map[string]MQTT.Client{}

	for name, connection := range parsedConfig.Connections.MQTT {
		result[name] = buildMqttConnection(connection)
	}

	return result
}

func buildMqttConnection(connection config.MQTTConnection) MQTT.Client {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(connection.Broker)

	opts.SetAutoReconnect(true)
	opts.SetClientID(connection.ClientId)
	opts.SetUsername(connection.Username)
	opts.SetPassword(connection.Password)
	opts.SetOnConnectHandler(onMqttConnect)
	opts.SetConnectionLostHandler(onMqttConnectionLost)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		zap.L().Fatal("Error while connecting to mqtt broker: %s", zap.Error(token.Error()))
	}

	return client
}

func onMqttConnect(client MQTT.Client) {
	zap.L().Info("MQTT connection established broker.")
}

func onMqttConnectionLost(client MQTT.Client, err error) {
	zap.L().Warn("Connection lost to broker.", zap.Error(err))
}
