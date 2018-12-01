package mqtt

import (
	"github.com/e154/smart-home-node/system/config"
)

type MqttConfig struct {
	SrvKeepAlive        int
	SrvConnectTimeout   int
	SrvSessionsProvider string
	SrvAuthenticator    string
	SrvTopicsProvider   string
	SrvPort             int
	SrvIp               string
}

func NewMqttConfig(cfg *config.AppConfig) *MqttConfig {
	return &MqttConfig{
		SrvKeepAlive:        cfg.MqttKeepAlive,
		SrvConnectTimeout:   cfg.MqttConnectTimeout,
		SrvSessionsProvider: cfg.MqttSessionsProvider,
		SrvAuthenticator:    cfg.MqttAuthenticator,
		SrvTopicsProvider:   cfg.MqttTopicsProvider,
		SrvPort:             cfg.MqttPort,
		SrvIp:               cfg.MqttIp,
	}
}
