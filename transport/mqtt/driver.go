package mqtt

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/arc/arc"
	arc_config "gitHub.***REMOVED***/monsoon/arc/config"
)

type MQTTClient struct {
	client       *MQTT.Client
	identity     string
	project      string
	organization string
	connected    bool
}

func New(config arc_config.Config) (*MQTTClient, error) {
	stdLogger := logrus.StandardLogger()
	logger := logrus.New()
	logger.Out = stdLogger.Out
	logger.Formatter = stdLogger.Formatter
	logger.Level = logrus.InfoLevel
	// We should really close this writer at some point
	w := logger.Writer()
	if logrus.GetLevel() >= logrus.ErrorLevel {
		MQTT.CRITICAL = log.New(w, "MQTT CRITICAL ", 0)
	}
	if logrus.GetLevel() >= logrus.ErrorLevel {
		MQTT.ERROR = log.New(w, "MQTT ERROR ", 0)
	}
	if logrus.GetLevel() >= logrus.DebugLevel {
		MQTT.WARN = log.New(w, "MQTT INFO ", 0)
	}
	if logrus.GetLevel() >= logrus.DebugLevel {
		MQTT.DEBUG = log.New(w, "MQTT DEBUG ", 0)
	}

	// set option
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
	if reg, err := offlineMessage(config.Organization, config.Project, config.Identity); err == nil {
		if j, err := reg.ToJSON(); err == nil {
			logrus.Infof("Setting last will delivering to %s", registrationTopic(config.Organization, config.Project, config.Identity))
			opts.SetBinaryWill(registrationTopic(config.Organization, config.Project, config.Identity), j, 0, false)
		}
	}
	opts.SetCleanSession(true)

	// create own transport
	transport := &MQTTClient{identity: config.Identity, project: config.Project, organization: config.Organization, connected: false}

	// set callbacks
	opts.OnConnect = func(_ *MQTT.Client) {
		transport.onConnect()
	}
	opts.OnConnectionLost = func(_ *MQTT.Client, err error) {
		transport.onConnectionLost(err)
	}

	// create client
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
	if reg, err := offlineMessage(c.organization, c.project, c.identity); err == nil {
		logrus.Info("Sending offline message")
		c.Registration(reg)
	} else {
		logrus.Error("Failed to create 'offline' registration message: ", err)
	}
	c.client.Disconnect(1000)
}

func (c *MQTTClient) IsConnected() bool {
	return c.connected && c.client.IsConnected()
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

func (c *MQTTClient) Request(msg *arc.Request) error {
	topic := identityTopic(msg.To)
	j, err := msg.ToJSON()
	if err != nil {
		logrus.Errorf("Error serializing Request to JSON: %s", err)
		return err
	} else {
		logrus.Debugf("Publishing request for %s/%s to %s", msg.Agent, msg.Action, topic)
		c.client.Publish(topic, 0, false, j).WaitTimeout(500 * time.Millisecond)
	}
	return nil
}

func (c *MQTTClient) Reply(msg *arc.Reply) error {
	topic := replyTopic(msg.RequestID)
	j, err := msg.ToJSON()
	if err != nil {
		logrus.Errorf("Error serializing Reply to JSON: %s", err)
		return err
	} else {
		logrus.Debugf("Publishing reply %s\n to %s", msg, topic)
		c.client.Publish(topic, 0, false, j).WaitTimeout(500 * time.Millisecond)
	}
	return nil
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

func (c *MQTTClient) Registration(msg *arc.Registration) error {
	topic := registrationTopic(c.organization, c.project, c.identity)
	j, err := msg.ToJSON()
	if err != nil {
		logrus.Errorf("Error serializing Registration to JSON: %s", err)
		return err
	} else {
		logrus.Debugf("Publishing registration %s\n to %s", msg, topic)
		c.client.Publish(topic, 0, false, j).WaitTimeout(500 * time.Millisecond)
	}
	return nil
}

func (c *MQTTClient) SubscribeRegistrations() (<-chan *arc.Registration, func()) {
	topic := "registration/+/+/+"
	out := make(chan *arc.Registration)
	canceled := make(chan struct{})
	var mutex sync.Mutex

	messageCallback := func(mClient *MQTT.Client, mMessage MQTT.Message) {
		payload := mMessage.Payload()
		msg, err := arc.ParseRegistration(&payload)
		if err != nil {
			logrus.Warnf("Discarding invalid message on topic %s:%s", mMessage.Topic(), err)
			return
		}
		if c := strings.Split(mMessage.Topic(), "/"); msg.Organization != c[1] || msg.Project != c[2] || msg.Sender != c[3] {
			logrus.Warn("Discarding message with modified organization, project or sender")
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

// Callbacks

func (c *MQTTClient) onConnect() {
	c.connected = true
}

func (c *MQTTClient) onConnectionLost(err error) {
	// send online message
	if req, err := onlineMessage(c.organization, c.project, c.identity); err == nil {
		logrus.Info("Sending online Message")
		c.Registration(req)
	} else {
		logrus.Error("Failed to create 'online' registration message ", err)
	}
	// set private state to true
	c.connected = false
}

// private

func identityTopic(identity string) string {
	return fmt.Sprintf("identity/%s", identity)
}
func replyTopic(request_id string) string {
	return fmt.Sprintf("reply/%s", request_id)
}

func registrationTopic(organization, project, identity string) string {
	return fmt.Sprintf("registration/%s/%s/%s", organization, project, identity)
}

func offlineMessage(organization, project, identity string) (*arc.Registration, error) {
	return arc.CreateRegistration(organization, project, identity, `{"online": false}`)
}
func onlineMessage(organization, project, identity string) (*arc.Registration, error) {
	return arc.CreateRegistration(organization, project, identity, `{"online": true}`)
}
