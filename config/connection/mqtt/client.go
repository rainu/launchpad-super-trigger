package mqtt

import MQTT "github.com/eclipse/paho.mqtt.golang"

type mqttClient struct {
	client      MQTT.Client
	conObserver connectionObserver
}

func (m *mqttClient) IsConnected() bool {
	return m.client.IsConnected()
}

func (m *mqttClient) IsConnectionOpen() bool {
	return m.client.IsConnectionOpen()
}

func (m *mqttClient) Connect() MQTT.Token {
	return m.client.Connect()
}

func (m *mqttClient) Disconnect(quiesce uint) {
	m.client.Disconnect(quiesce)
}

func (m *mqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	return m.client.Publish(topic, qos, retained, payload)
}

func (m *mqttClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	return m.client.Subscribe(topic, qos, callback)
}

func (m *mqttClient) SubscribeMultiple(filters map[string]byte, callback MQTT.MessageHandler) MQTT.Token {
	return m.client.SubscribeMultiple(filters, callback)
}

func (m *mqttClient) Unsubscribe(topics ...string) MQTT.Token {
	return m.client.Unsubscribe(topics...)
}

func (m *mqttClient) AddRoute(topic string, callback MQTT.MessageHandler) {
	m.client.AddRoute(topic, callback)
}

func (m *mqttClient) OptionsReader() MQTT.ClientOptionsReader {
	return m.client.OptionsReader()
}
