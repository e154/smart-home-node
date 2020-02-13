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

package main

import (
	"github.com/e154/smart-home-node/system/client"
	"github.com/e154/smart-home-node/system/config"
	"github.com/e154/smart-home-node/system/dig"
	"github.com/e154/smart-home-node/system/graceful_service"
	"github.com/e154/smart-home-node/system/logging"
	"github.com/e154/smart-home-node/system/mqtt"
	"github.com/e154/smart-home-node/system/serial"
	"github.com/e154/smart-home-node/system/tcpproxy"
)

func BuildContainer() (container *dig.Container) {

	container = dig.New()
	container.Provide(logging.NewLogrus)
	container.Provide(config.ReadConfig)
	container.Provide(graceful_service.NewGracefulService)
	container.Provide(graceful_service.NewGracefulServicePool)
	container.Provide(graceful_service.NewGracefulServiceConfig)
	container.Provide(mqtt.NewMqtt)
	container.Provide(mqtt.NewMqttConfig)
	container.Provide(client.NewClient)
	container.Provide(serial.NewSerialService)
	container.Provide(tcpproxy.NewTcpProxy)

	return
}
