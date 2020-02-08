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

package config

type AppConfig struct {
	Name                 string   `json:"name"`
	MqttClientId         string   `json:"mqtt_client_id"`
	MqttKeepAlive        int      `json:"mqtt_keep_alive"`
	MqttConnectTimeout   int      `json:"mqtt_connect_timeout"`
	MqttSessionsProvider string   `json:"mqtt_sessions_provider"`
	MqttUsername         string   `json:"mqtt_username"`
	MqttPassword         string   `json:"mqtt_password"`
	MqttTopicsProvider   string   `json:"mqtt_topics_provider"`
	MqttIp               string   `json:"mqtt_ip"`
	MqttPort             int      `json:"mqtt_port"`
	ProxyPort            int      `json:"proxy_port"`
	Serial               []string `json:"serial"`
}

type RunMode string

const (
	DebugMode   = RunMode("debug")
	ReleaseMode = RunMode("release")
)
