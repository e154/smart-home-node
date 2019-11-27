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
	}
}
