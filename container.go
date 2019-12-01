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
