package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

type Client struct {
	qos            byte
	baseTopic, uri string
	client         MQTT.Client
	cfg            *MqttConfig
	handler        func(MQTT.Client, MQTT.Message)
}

func NewClient(cfg *MqttConfig, baseTopic string, qos byte,
	handler func(MQTT.Client, MQTT.Message)) (client *Client, err error) {

	uri := fmt.Sprintf("tcp://%s:%d", cfg.ServerIp, cfg.Port)

	client = &Client{
		baseTopic: baseTopic,
		handler:   handler,
		qos:       qos,
		uri:       uri,
		cfg:       cfg,
	}

	clientId := fmt.Sprintf("node-%d-%d", os.Getpid(), time.Now().Unix())
	opts := MQTT.NewClientOptions().
		AddBroker(uri).
		SetClientID(clientId).
		SetKeepAlive(time.Duration(cfg.KeepAlive) * time.Second).
		SetPingTimeout(5 * time.Second).
		SetConnectTimeout(time.Duration(cfg.ConnectTimeout) * time.Second).
		SetCleanSession(true).
		SetUsername(cfg.MqttUsername).
		SetPassword(cfg.MqttPassword).
		SetOnConnectHandler(client.onConnect).
		SetConnectionLostHandler(client.onConnectionLostHandler)

	client.client = MQTT.NewClient(opts)

	return
}

func (c *Client) onConnectionLostHandler(client MQTT.Client, e error) {

	log.Debug("connection lost...")

	c.unsubscribe()
}

func (c *Client) onConnect(client MQTT.Client) {

	log.Debug("connected...")

	if token := c.client.Subscribe(c.baseTopic+"/req", c.qos, c.handler); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
	}
}

func (c *Client) Connect() {

	log.Infof("Connect to server %s", c.uri)

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
	}

	return
}

func (c *Client) Disconnect() {
	if c.client == nil {
		return
	}

	c.unsubscribe()
	c.client.Disconnect(250)
	c.client = nil
}

func (c *Client) unsubscribe() {
	if token := c.client.Unsubscribe(c.baseTopic + "/req"); token.Error() != nil {
		log.Error(token.Error().Error())
	}
}

func (c *Client) Publish(topic string, payload interface{}) (err error) {
	if c.client != nil && (c.client.IsConnected()) {
		c.client.Publish(c.baseTopic+topic, c.qos, false, payload)
	}
	return
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnectionOpen()
}
