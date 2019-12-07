package mqtt

import (
	"fmt"
	"github.com/e154/smart-home-node/system/graceful_service"
	"github.com/e154/smart-home-node/system/mqtt_client"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("mqtt")
)

type Mqtt struct {
	cfg     *MqttConfig
	clients []*mqtt_client.Client
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

func (m *Mqtt) NewClient(cfg *mqtt_client.Config) (c *mqtt_client.Client, err error) {

	if cfg == nil {
		cfg = &mqtt_client.Config{
			KeepAlive:      m.cfg.KeepAlive,
			PingTimeout:    5,
			Broker:         fmt.Sprintf("tcp://%s:%d", m.cfg.ServerIp, m.cfg.Port),
			ClientID:       m.cfg.MqttClientId,
			ConnectTimeout: m.cfg.ConnectTimeout,
			CleanSession:   true,
			Username:       m.cfg.MqttUsername,
			Password:       m.cfg.MqttPassword,
		}
	}

	if c, err = mqtt_client.NewClient(cfg); err != nil {
		return
	}

	m.clients = append(m.clients, c)

	return
}
