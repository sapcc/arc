package mqtt

import (
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	//logrus "github.com/Sirupsen/logrus"
	"gitHub.***REMOVED***/monsoon/onos/types"
	"log"
	"os"
)

type MQTTClient struct {
	client *MQTT.Client
	identity string
	project string
}

func New(config types.Config) (*MQTTClient, error) {

	MQTT.CRITICAL = log.New(os.Stdout, "MQTT CRITICAL", log.LstdFlags)
	MQTT.ERROR = log.New(os.Stdout, "MQTT ERROR", log.LstdFlags)
	MQTT.WARN = log.New(os.Stdout, "MQTT INFO", log.LstdFlags)
	MQTT.DEBUG = log.New(os.Stdout, "MQTT DEBUG", log.LstdFlags)

	opts := MQTT.NewClientOptions()
	for _, endpoint := range config.Endpoints {
		opts.AddBroker(endpoint)
	}
	opts.SetCleanSession(true)
	c := MQTT.NewClient(opts)
	return &MQTTClient{client: c, identity: config.Identity, project: config.Project}, nil
}

func (c *MQTTClient) Connect() {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
	}
}

func (c *MQTTClient) Disconnect() {
	c.Disconnect()
}

func (c *MQTTClient) Subscribe( callback func(types.Message)) {
	c.client.Subscribe('test', 
}

func (c *MQTTClient) Publish(msg types.Message) {
}
