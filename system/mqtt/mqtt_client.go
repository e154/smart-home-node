package mqtt

import (
	"fmt"
	"time"
	"os"
	"github.com/surgemq/surgemq/service"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"errors"
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
		SetCleanSession(true)

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
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
		time.Sleep(time.Second)
		goto loop
		return
	}

	//time.Sleep(time.Second)
	if err = c.Subscribe(c.topic+"/req", c.qos, c.handler); err != nil {
		log.Warning(err.Error())
		time.Sleep(time.Second)
		goto loop
	}

	//if token := c.client.Subscribe("$SYS/broker/connection/#", 0, c.brokerConnectionHandler); token.Wait() && token.Error() != nil {
	//	log.Error(token.Error().Error())
	//	time.Sleep(time.Second)
	//	c.Connect()
	//	return
	//}
	//if token := c.client.Subscribe("$SYS/broker/load/#", 0, c.brokerLoadHandler); token.Wait() && token.Error() != nil {
	//	log.Error(token.Error().Error())
	//	time.Sleep(time.Second)
	//	c.Connect()
	//	return
	//}
	//if token := c.client.Subscribe("$SYS/broker/clients/#", 0, c.brokerClientsHandler); token.Wait() && token.Error() != nil {
	//	log.Error(token.Error().Error())
	//	time.Sleep(time.Second)
	//	c.Connect()
	//	return
	//}

	return
}

func (c *Client) Disconnect() {
	if c.client != nil && (c.client.IsConnected()) {
		c.client.Disconnect(250)
	}
}

func (c *Client) Publish(topic string, payload interface{}) (err error) {
	if c.client != nil && (c.client.IsConnected()) {
		c.client.Publish(c.topic + topic, c.qos, false, payload)
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

//func (c *Client) brokerLoadHandler(client MQTT.Client, msg MQTT.Message) {
//	c.brokerLoad <- true
//	log.Debugf("BrokerLoadHandler         ")
//	log.Debugf("[%s]  ", msg.Topic())
//	log.Debugf("%s\n", msg.Payload())
//}
//func (c *Client) brokerConnectionHandler(client MQTT.Client, msg MQTT.Message) {
//	c.brokerConnection <- true
//	log.Debugf("BrokerConnectionHandler   ")
//	log.Debugf("[%s]  ", msg.Topic())
//	log.Debugf("%s\n", msg.Payload())
//}
//func (c *Client) brokerClientsHandler(client MQTT.Client, msg MQTT.Message) {
//	c.brokerClients <- true
//	log.Debugf("BrokerClientsHandler      ")
//	log.Debugf("[%s]  ", msg.Topic())
//	log.Debugf("%s\n", msg.Payload())
//}

func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}