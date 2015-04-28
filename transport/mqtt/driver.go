package mqtt

import (
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type MQTTClient struct {
	client *MQTT.Client
}

func New() (*MQTTClient, error) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	c := MQTT.NewClient(opts)
	return &MQTTClient{c}, nil
}

func (c *MQTTClient) Connect() {
}

func (c *MQTTClient) Disconnect() {
}
