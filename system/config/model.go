package config

type AppConfig struct {
	Name                 string   `json:"name"`
	Topic                string   `json:"topic"`
	MqttKeepAlive        int      `json:"mqtt_keep_alive"`
	MqttConnectTimeout   int      `json:"mqtt_connect_timeout"`
	MqttSessionsProvider string   `json:"mqtt_sessions_provider"`
	MqttAuthenticator    string   `json:"mqtt_authenticator"`
	MqttTopicsProvider   string   `json:"mqtt_topics_provider"`
	MqttIp               string   `json:"mqtt_ip"`
	MqttPort             int      `json:"mqtt_port"`
	Serial               []string `json:"serial"`
}

type RunMode string

const (
	DebugMode   = RunMode("debug")
	ReleaseMode = RunMode("release")
)
