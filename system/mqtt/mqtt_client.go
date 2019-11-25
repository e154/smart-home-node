package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

type Client struct {
	qos              byte
	baseTopic, uri   string
	clientId         string
	client           MQTT.Client
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

func NewClient(cfg *MqttConfig, baseTopic string, qos byte,
	handler func(MQTT.Client, MQTT.Message)) (client *Client, err error) {

	// Instantiates a new Client
	uri := fmt.Sprintf("tcp://%s:%d", cfg.ServerIp, cfg.Port)

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
		baseTopic:        baseTopic,
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

//loop:
	time.Sleep(time.Millisecond * 500)

	log.Info("connect ....")

	//if c.client == nil {
	//	return
	//}

	if c.client.IsConnected() {
		c.client.Disconnect(250)
	}

	if c.client.IsConnected() {
		return
	}

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
		//goto loop
	}

	if token := c.Subscribe(c.baseTopic+"/req", c.qos, c.handler); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
		//goto loop
	}

	if token := c.Subscribe(c.baseTopic+"/pong", c.qos, c.pong); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
		//goto loop
	}

	c.lastPing = time.Now()
	c.reconnect = false

	return
}

func (c *Client) Disconnect() {
	if c.client == nil {
		return
	}

	c.quit <- struct{}{}

	if token := c.client.Unsubscribe(c.baseTopic + "/req"); token.Error() != nil {
		log.Error(token.Error().Error())
	}

	if token := c.client.Unsubscribe(c.baseTopic + "/pong"); token.Error() != nil {
		log.Error(token.Error().Error())
	}

	c.client.Disconnect(250)

	c.client = nil
}

func (c *Client) Publish(topic string, payload interface{}) (err error) {
	if c.client != nil && (c.client.IsConnected()) {
		c.client.Publish(c.baseTopic+topic, c.qos, false, payload)
	}
	return
}

func (c *Client) Subscribe(topic string, qos byte, handler func(MQTT.Client, MQTT.Message)) MQTT.Token {
	return c.client.Subscribe(topic, qos, handler)
}

func (c *Client) IsConnected() bool {
	//if c.reconnect {
	//	return false
	//}

	k := time.Now().Sub(c.lastPing).Seconds()
	if ok := k < 5; !ok {
		//c.reconnect = true
		//c.Connect()
		return false
	}

	return true
}

func (c *Client) pong(MQTT.Client, MQTT.Message) {
	c.lastPing = time.Now()
}
