// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package mqtt

import (
	"fmt"
	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/system/graceful_service"
	"github.com/e154/smart-home-node/system/mqtt_client"
)

var (
	log = common.MustGetLogger("mqtt")
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
