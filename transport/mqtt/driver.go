package mqtt

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/arc/arc"
)

type MQTTClient struct {
	client       *MQTT.Client
	identity     string
	project      string
	organization string
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
	if len(config.Endpoints) == 0 {
		return nil, fmt.Errorf("No transport endpoints given")
	}
	//check the first endpoint if we need to setup a TlsConfig
	if url, err := url.Parse(config.Endpoints[0]); err == nil {
		switch url.Scheme {
		case "tls", "ssl", "tcps", "wss":

			tlsc := tls.Config{
				RootCAs:      config.CACerts,
				Certificates: []tls.Certificate{*config.ClientCert},
				MinVersion:   tls.VersionTLS12,
			}
			opts.SetTLSConfig(&tlsc)
		}

	} else {
		return nil, fmt.Errorf("Invalid url as transport endpoint given")
	}
	for _, endpoint := range config.Endpoints {
		logrus.Info("Using MQTT broker ", endpoint)
		opts.AddBroker(endpoint)
	}
	if reg, err := offlineMessage(config.Identity); err == nil {
		if j, err := reg.ToJSON(); err == nil {
			logrus.Infof("Setting last will delivering to %s", identityTopic(reg.To))
			opts.SetBinaryWill(identityTopic(reg.To), j, 0, false)
		}
	}
	transport := &MQTTClient{identity: config.Identity, project: config.Project, organization: config.Organization}
	if req, err := onlineMessage(config.Identity); err == nil {
		opts.OnConnect = func(_ *MQTT.Client) {
			logrus.Info("Sending online Message")
			transport.Request(req)
		}
	} else {
		logrus.Error("Failed to create 'online' registration message ", err)

	}
	opts.SetCleanSession(true)
	c := MQTT.NewClient(opts)
	transport.client = c

	return transport, nil
}

func (c *MQTTClient) Connect() error {
	logrus.Info("Connecting to MQTT broker")
	token := c.client.Connect()
	token.Wait()
	return token.Error()
}

func (c *MQTTClient) Disconnect() {
	if req, err := offlineMessage(c.identity); err == nil {
		logrus.Info("Sending offline message")
		c.Request(req)
	} else {
		logrus.Error("Failed to create 'offline' registration message: ", err)
	}
	c.client.Disconnect(1000)
}

func (c *MQTTClient) Subscribe(identity string) (<-chan *arc.Request, func()) {
	topic := identityTopic(identity)
	msgChan := make(chan *arc.Request)
	canceled := make(chan struct{})
	var mutex sync.Mutex

	messageCallback := func(mClient *MQTT.Client, mMessage MQTT.Message) {
		payload := mMessage.Payload()
		msg, err := arc.ParseRequest(&payload)
		if err != nil {
			logrus.Warnf("Discarding invalid message on topic %s:%s", mMessage.Topic(), err)
			return
		}
		mutex.Lock()
		select {
		case <-canceled:
		case msgChan <- msg:
		}
		mutex.Unlock()
	}

	cancel := func() {
		c.client.Unsubscribe(topic).Wait()
		close(canceled)
		mutex.Lock()
		close(msgChan)
		mutex.Unlock()
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
		logrus.Debugf("Publishing request for %s/%s to %s", msg.Agent, msg.Action, topic)
		c.client.Publish(topic, 0, false, j).WaitTimeout(500 * time.Millisecond)
	}
}

func (c *MQTTClient) Reply(msg *arc.Reply) {
	topic := replyTopic(msg.RequestID)
	j, err := msg.ToJSON()
	if err != nil {
		logrus.Errorf("Error serializing Reply to JSON: %s", err)
	} else {
		logrus.Debugf("Publishing reply %s\n to %s", msg, topic)
		c.client.Publish(topic, 0, false, j).WaitTimeout(500 * time.Millisecond)
	}
}

func (c *MQTTClient) SubscribeJob(requestId string) (<-chan *arc.Reply, func()) {
	topic := replyTopic(requestId)
	out := make(chan *arc.Reply)
	canceled := make(chan struct{})
	var mutex sync.Mutex

	messageCallback := func(mClient *MQTT.Client, mMessage MQTT.Message) {
		payload := mMessage.Payload()
		msg, err := arc.ParseReply(&payload)
		if err != nil {
			logrus.Warnf("Discarding invalid message on topic %s:%s", mMessage.Topic(), err)
			return
		}
		mutex.Lock()
		select {
		case <-canceled:
		case out <- msg:
		}
		mutex.Unlock()
	}
	cancel := func() {
		c.client.Unsubscribe(topic).Wait()
		close(canceled)
		mutex.Lock()
		close(out)
		mutex.Unlock()
	}

	c.client.Subscribe(topic, 0, messageCallback).Wait()
	return out, cancel
}

func (c *MQTTClient) SubscribeReplies() (<-chan *arc.Reply, func()) {
	//This is a little bit hacky but YOLO
	//At some point we need to rethink that, maybe have "namespaces" one can subscribe to
	return c.SubscribeJob("+")
}

func identityTopic(identity string) string {
	return fmt.Sprintf("identity/%s", identity)
}
func replyTopic(request_id string) string {
	return fmt.Sprintf("reply/%s", request_id)
}

func offlineMessage(identity string) (*arc.Request, error) {
	return arc.CreateRegistrationMessage(identity, `{"online": false}`)
}
func onlineMessage(identity string) (*arc.Request, error) {
	return arc.CreateRegistrationMessage(identity, `{"online": true}`)
}
