package mqtt

import (
	"encoding/json"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/onos/onos"
	"log"
)

type MQTTClient struct {
	client   *MQTT.Client
	identity string
	project  string
}

func New(config onos.Config) (*MQTTClient, error) {
	stdLogger := logrus.StandardLogger()
	logger := logrus.New()
	logger.Out = stdLogger.Out
	logger.Formatter = stdLogger.Formatter
	logger.Level = logrus.InfoLevel
	// We should really close this writer at some point
	w := logger.Writer()
	if logrus.GetLevel() >= logrus.FatalLevel {
		MQTT.CRITICAL = log.New(w, "MQTT CRITICAL ", 0)
	}
	if logrus.GetLevel() >= logrus.ErrorLevel {
		MQTT.ERROR = log.New(w, "MQTT ERROR ", 0)
	}
	if logrus.GetLevel() >= logrus.InfoLevel {
		MQTT.WARN = log.New(w, "MQTT INFO ", 0)
	}
	if logrus.GetLevel() >= logrus.DebugLevel {
		MQTT.DEBUG = log.New(w, "MQTT DEBUG ", 0)
	}

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

func (c *MQTTClient) Subscribe() <-chan *onos.Message {
	msgChan := make(chan *onos.Message)
	messageCallback := func(mClient *MQTT.Client, mMessage MQTT.Message) {
		msg, err := parseMessage(mMessage)
		if err != nil {
			logrus.Warnf("Discarding invalid message on topic %s:%s\n", mMessage.Topic(), err)
			return
		}
		msgChan <- &msg
	}
	c.client.Subscribe("test", 0, messageCallback)
	return msgChan
}

func (c *MQTTClient) Publish(msg onos.Message) {
}

// private

func parseMessage(msg MQTT.Message) (onos.Message, error) {
	var m onos.Message
	err := json.Unmarshal(msg.Payload(), &m)
	return m, err
}
