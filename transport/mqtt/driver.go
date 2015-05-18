package mqtt

import (
	"fmt"
	"log"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type MQTTClient struct {
	client   *MQTT.Client
	identity string
	project  string
}

func New(config arc.Config) (*MQTTClient, error) {
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

func (c *MQTTClient) Subscribe() <-chan *arc.Request {
	msgChan := make(chan *arc.Request)
	messageCallback := func(mClient *MQTT.Client, mMessage MQTT.Message) {
		payload := mMessage.Payload()
		msg, err := arc.ParseRequest(&payload)
		if err != nil {
			logrus.Warnf("Discarding invalid message on topic %s:%s\n", mMessage.Topic(), err)
			return
		}
		msgChan <- msg
	}
	c.client.Subscribe("test", 0, messageCallback)
	return msgChan
}

func (c *MQTTClient) Request(msg *arc.Request) {
	logrus.Debug("Publishing request %s\n", msg)
}

func (c *MQTTClient) Reply(msg *arc.Reply) {
	var topic = fmt.Sprintf("reply/%s", msg.RequestID)
	logrus.Debugf("Publishing reply %s\n to %s", msg, topic)
	j, err := msg.ToJSON()
	if err != nil {
		logrus.Errorf("Error serializing Reply to JSON: %s", err)
	} else {
		c.client.Publish(topic, 0, false, j)
	}
}

func (c *MQTTClient) SubscribeJob(requestId string) <-chan *arc.Reply {
	return make(chan *arc.Reply)
}
