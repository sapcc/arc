package mqtt

import (
	"encoding/json"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/onos/onos"
	"log"
	"os"
)

type MQTTClient struct {
	client   *MQTT.Client
	identity string
	project  string
}

func New(config onos.Config) (*MQTTClient, error) {

	MQTT.CRITICAL = log.New(os.Stdout, "MQTT CRITICAL", log.LstdFlags)
	MQTT.ERROR = log.New(os.Stdout, "MQTT ERROR", log.LstdFlags)
	MQTT.WARN = log.New(os.Stdout, "MQTT INFO", log.LstdFlags)
	MQTT.DEBUG = log.New(os.Stdout, "MQTT DEBUG", log.LstdFlags)

	opts := MQTT.NewClientOptions()
	for _, endpoint := range config.Endpoints {
		logrus.Info("Using MQTT broker ", endpoint)
		opts.AddBroker(endpoint)
	}
	opts.SetCleanSession(true)
	c := MQTT.NewClient(opts)
	return &MQTTClient{client: c, identity: config.Identity, project: config.Project}, nil
}

func (c *MQTTClient) Connect() {
	logrus.Info("Connecting to MQTT broker")
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
	}
}

func (c *MQTTClient) Disconnect() {
	c.client.Disconnect(1000)
}

func (c *MQTTClient) Subscribe(callback func(onos.Message)) {

	messageHandler := func(mClient *MQTT.Client, mMessage MQTT.Message) {
		msg, err := parseMessage(mMessage)
		if err != nil {
			logrus.Warnf("Discarding invalid message on topic %s:%s\n", mMessage.Topic(), err)
			return
		}
		logrus.Infof("Received message with requestId %s\n", msg.RequestId)
		callback(msg)
	}
	c.client.Subscribe("test", 0, messageHandler)
}

func (c *MQTTClient) Publish(msg onos.Message) {
}

// private

func parseMessage(msg MQTT.Message) (onos.Message, error) {
	var m onos.Message
	err := json.Unmarshal(msg.Payload(), &m)
	return m, err
}
