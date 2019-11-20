package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/surgemq/surgemq/service"
	"os"
	"time"
)

type Client struct {
	qos              byte
	topic, uri       string
	client           MQTT.Client
	onComplete       service.OnCompleteFunc
	onPublish        service.OnPublishFunc
	brokerLoad       chan bool
	brokerConnection chan bool
	brokerClients    chan bool
	pongChannel      chan struct{}
	cfg              *MqttConfig
	loadCount        int
	connectionCount  int
	clientsCount     int
	handler          func(MQTT.Client, MQTT.Message)
	lastPing         time.Time
	reconnect        bool
	quit             chan struct{}
}

func NewClient(cfg *MqttConfig,
	topic string,
	qos byte,
	handler func(MQTT.Client, MQTT.Message)) (client *Client, err error) {

	// Instantiates a new Client
	uri := fmt.Sprintf("tcp://%s:%d", cfg.SrvIp, cfg.SrvPort)

	clientId := fmt.Sprintf("node-%d-%d", os.Getpid(), time.Now().Unix())
	opts := MQTT.NewClientOptions().
		AddBroker(uri).
		SetClientID(clientId).
		SetKeepAlive(2 * time.Second).
		SetPingTimeout(1 * time.Second).
		SetCleanSession(true).
		SetUsername(cfg.MqttUsername).
		SetPassword(cfg.MqttPassword)

	c := MQTT.NewClient(opts)

	client = &Client{
		topic:            topic,
		handler:          handler,
		qos:              qos,
		client:           c,
		uri:              uri,
		brokerLoad:       make(chan bool),
		brokerConnection: make(chan bool),
		brokerClients:    make(chan bool),
		pongChannel:      make(chan struct{}),
		cfg:              cfg,
		reconnect:        true,
		quit:             make(chan struct{}),
	}

	go func() {
		for ; ; {
			select {
			case <-client.brokerLoad:
				client.loadCount++
			case <-client.brokerConnection:
				client.connectionCount++
			case <-client.brokerClients:
				client.clientsCount++
			case <-client.quit:
				return
			}
		}
	}()

	return
}

func (c *Client) Connect() {

	log.Infof("Connect to server %s", c.uri)

loop:
	time.Sleep(time.Second)

	log.Info("connect ....")

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		goto loop
	}

	log.Info("connect ....1")

	if token := c.Subscribe(c.topic+"/req", c.qos, c.handler); token.Error() != nil {
		goto loop
	}

	log.Info("connect ....2")

	if token := c.Subscribe(c.topic+"/pong", c.qos, c.pong); token.Error() != nil {
		goto loop
	}

	log.Info("connect ....3")

	c.lastPing = time.Now()
	c.reconnect = false

	return
}

func (c *Client) Disconnect() {
	if c.client == nil {
		return
	}

	c.quit <- struct{}{}

	if token := c.client.Unsubscribe(c.topic + "/req"); token.Error() != nil {
		log.Error(token.Error().Error())
	}

	if token := c.client.Unsubscribe(c.topic + "/pong"); token.Error() != nil {
		log.Error(token.Error().Error())
	}

	c.client.Disconnect(250)
}

func (c *Client) Publish(topic string, payload interface{}) (err error) {
	if c.client != nil && (c.client.IsConnected()) {
		c.client.Publish(c.topic+topic, c.qos, false, payload)
	}
	return
}

func (c *Client) Subscribe(topic string, qos byte, handler func(MQTT.Client, MQTT.Message)) MQTT.Token {
	return c.client.Subscribe(topic, qos, handler)
}

func (c *Client) IsConnected() bool {
	if c.reconnect {
		return false
	}

	k := time.Now().Sub(c.lastPing).Seconds()
	if ok := k < 5; !ok {
		c.reconnect = true
		c.Connect()
		return false
	}

	return true
}

func (c *Client) pong(MQTT.Client, MQTT.Message) {
	c.lastPing = time.Now()
}
