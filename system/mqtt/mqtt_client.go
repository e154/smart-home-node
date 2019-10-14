package mqtt

import (
	"errors"
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
	cfg              *MqttConfig
	loadCount        int
	connectionCount  int
	clientsCount     int
	handler          func(MQTT.Client, MQTT.Message)
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
		cfg:              cfg,
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
			}
		}
	}()

	//go client.ping()

	return
}

func (c *Client) Connect() (err error) {

	log.Infof("Connect to server %s", c.uri)

loop:
	time.Sleep(time.Second)

	log.Info("connect ....")

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		goto loop
	}

	log.Info("connect ....1")

	if err = c.Subscribe(c.topic+"/req", c.qos, c.handler); err != nil {
		log.Warning(err.Error())
		goto loop
	}

	log.Info("connect ....2")

	if err = c.Subscribe("$SYS/broker/connection/#", 0, func(client MQTT.Client, message MQTT.Message) {

	}); err != nil {
		log.Warning(err.Error())
		goto loop
	}

	return
}

func (c *Client) Disconnect() {
	if c.client != nil && (c.client.IsConnected()) {
		c.client.Disconnect(250)
	}
}

func (c *Client) Publish(topic string, payload interface{}) (err error) {
	if c.client != nil && (c.client.IsConnected()) {
		c.client.Publish(c.topic+topic, c.qos, false, payload)
	}
	return
}

func (c *Client) Subscribe(topic string, qos byte, handler func(MQTT.Client, MQTT.Message)) (err error) {

	if token := c.client.Subscribe(topic, qos, handler); token.Wait() && token.Error() != nil {
		err = errors.New(token.Error().Error())
		return
	}

	return
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}
