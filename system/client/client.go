package client

import (
	"fmt"
	"github.com/e154/smart-home-node/system/config"
	"github.com/e154/smart-home-node/system/graceful_service"
	"github.com/op/go-logging"
	"github.com/e154/smart-home-node/system/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"encoding/json"
	"github.com/e154/smart-home-node/common/debug"
)

var (
	log = logging.MustGetLogger("client")
)

type Client struct {
	cfg    *config.AppConfig
	client *mqtt.Client
}

func NewClient(cfg *config.AppConfig,
	graceful *graceful_service.GracefulService,
	qService *mqtt.Mqtt) *Client {

	client := &Client{
		cfg:    cfg,
	}
	topic := fmt.Sprintf("/home/%s", cfg.Topic)
	c, _ := qService.NewClient(topic, 0x0, client.onPublish)
	client.client = c

	graceful.Subscribe(client)

	return client
}

func (c *Client) Shutdown() {
	c.client.Disconnect()
}

func (c *Client) Connect() {
	go c.client.Connect()
}

func (c *Client) onPublish(cli MQTT.Client, msg MQTT.Message) {

	message := &Message{}
	if err := json.Unmarshal(msg.Payload(), message); err != nil {
		log.Error(err.Error())
	}

	debug.Println(message)
}
