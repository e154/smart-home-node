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
	"github.com/e154/smart-home-node/system/config"
)

type MqttConfig struct {
	KeepAlive        int
	ConnectTimeout   int
	SessionsProvider string
	MqttUsername     string
	MqttPassword     string
	TopicsProvider   string
	Port             int
	ServerIp         string
	MqttClientId     string
}

func NewMqttConfig(cfg *config.AppConfig) *MqttConfig {
	return &MqttConfig{
		KeepAlive:        cfg.MqttKeepAlive,
		ConnectTimeout:   cfg.MqttConnectTimeout,
		SessionsProvider: cfg.MqttSessionsProvider,
		MqttUsername:     cfg.MqttUsername,
		MqttPassword:     cfg.MqttPassword,
		TopicsProvider:   cfg.MqttTopicsProvider,
		Port:             cfg.MqttPort,
		ServerIp:         cfg.MqttIp,
		MqttClientId:     cfg.MqttClientId,
	}
}
