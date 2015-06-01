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
	if logrus.GetLevel() >= logrus.DebugLevel {
		MQTT.CRITICAL = log.New(w, "MQTT CRITICAL ", 0)
	}
	if logrus.GetLevel() >= logrus.DebugLevel {
		MQTT.ERROR = log.New(w, "MQTT ERROR ", 0)
	}
	if logrus.GetLevel() >= logrus.DebugLevel {
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

func (c *MQTTClient) Connect() error {
	logrus.Info("Connecting to MQTT broker")
	token := c.client.Connect()
	token.Wait()
	return token.Error()
}

func (c *MQTTClient) Disconnect() {
	c.client.Disconnect(1000)
}

func (c *MQTTClient) Subscribe(identity string) (<-chan *arc.Request, func()) {
	topic := identityTopic(identity)
	msgChan := make(chan *arc.Request)
	messageCallback := func(mClient *MQTT.Client, mMessage MQTT.Message) {
		payload := mMessage.Payload()
		msg, err := arc.ParseRequest(&payload)
		if err != nil {
			logrus.Warnf("Discarding invalid message on topic %s:%s", mMessage.Topic(), err)
			return
		}
		msgChan <- msg
	}

	cancel := func() {
		c.client.Unsubscribe(topic).Wait()
		close(msgChan)
	}

	c.client.Subscribe(topic, 0, messageCallback).Wait()
	return msgChan, cancel
}

func (c *MQTTClient) Request(msg *arc.Request) {
	topic := identityTopic(msg.To)
	j, err := msg.ToJSON()
	if err != nil {
		logrus.Errorf("Error serializing Request to JSON: %s", err)
	} else {
		logrus.Debug("Publishing request %s\n", msg)
		c.client.Publish(topic, 0, false, j)
	}
}

func (c *MQTTClient) Reply(msg *arc.Reply) {
	topic := replyTopic(msg.RequestID)
	j, err := msg.ToJSON()
	if err != nil {
		logrus.Errorf("Error serializing Reply to JSON: %s", err)
	} else {
		logrus.Debugf("Publishing reply %s\n to %s", msg, topic)
		c.client.Publish(topic, 0, false, j)
	}
}

func (c *MQTTClient) SubscribeJob(requestId string) (<-chan *arc.Reply, func()) {
	topic := replyTopic(requestId)
	out := make(chan *arc.Reply)

	messageCallback := func(mClient *MQTT.Client, mMessage MQTT.Message) {
		payload := mMessage.Payload()
		msg, err := arc.ParseReply(&payload)
		if err != nil {
			logrus.Warnf("Discarding invalid message on topic %s:%s", mMessage.Topic(), err)
			return
		}
		out <- msg
	}
	cancel := func() {
		c.client.Unsubscribe(topic).Wait()
		close(out)
	}

	c.client.Subscribe(topic, 0, messageCallback).Wait()
	return out, cancel
}

func identityTopic(identity string) string {
	return fmt.Sprintf("identity/%s", identity)
}
func replyTopic(request_id string) string {
	return fmt.Sprintf("reply/%s", request_id)
}
