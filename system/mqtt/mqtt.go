package mqtt

import (
	"github.com/e154/smart-home-node/system/graceful_service"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("mqtt")
)

type Mqtt struct {
	cfg     *MqttConfig
	clients []*Client
}

func NewMqtt(cfg *MqttConfig,
	graceful *graceful_service.GracefulService) (mqtt *Mqtt) {
	mqtt = &Mqtt{
		cfg: cfg,
	}

	graceful.Subscribe(mqtt)

	return
}

func (m *Mqtt) Shutdown() {
	for _, client := range m.clients {
		if client == nil {
			continue
		}
		client.Disconnect()
	}

	log.Info("Server exiting")
}

func (m *Mqtt) NewClient(baseTopic string, qos byte, handler func(MQTT.Client, MQTT.Message)) (c *Client, err error) {

	if c, err = NewClient(m.cfg, baseTopic, qos, handler); err != nil {
		return
	}

	m.clients = append(m.clients, c)

	return
}
