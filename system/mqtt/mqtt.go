package mqtt

import (
	"github.com/op/go-logging"
	"github.com/surgemq/surgemq/service"
	"github.com/e154/smart-home-node/system/graceful_service"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	log = logging.MustGetLogger("mqtt")
)

type Mqtt struct {
	cfg     *MqttConfig
	server  *service.Server
	clients []*Client
}

func NewMqtt(cfg *MqttConfig,
	graceful *graceful_service.GracefulService) (mqtt *Mqtt) {
	mqtt = &Mqtt{
		cfg: cfg,
	}

	//go mqtt.runServer()

	graceful.Subscribe(mqtt)

	return
}

func (m *Mqtt) Shutdown() {
	//if m.server != nil {
	//	m.server.Close()
	//}

	for _, client := range m.clients {
		if client == nil {
			continue
		}
		client.Disconnect()
	}

	log.Info("Server exiting")
}

func (m *Mqtt) NewClient(topic string,
	qos byte,
	handler func(MQTT.Client, MQTT.Message)) (c *Client, err error) {

	if c, err = NewClient(m.cfg, topic, qos, handler); err != nil {
		return
	}

	m.clients = append(m.clients, c)

	return
}
