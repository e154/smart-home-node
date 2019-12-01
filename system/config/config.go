package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const path = "conf/config.json"

func ReadConfig() (conf *AppConfig, err error) {
	var file []byte
	file, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading config file")
		return
	}
	conf = &AppConfig{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal("Error: wrong format of config file")
		return
	}

	checkEnv(conf)

	return
}

func checkEnv(conf *AppConfig) {

	if name := os.Getenv("NAME"); name != "" {
		conf.Name = name
	}

	if topic := os.Getenv("TOPIC"); topic != "" {
		conf.Topic = topic
	}

	if mqttKeepAlive := os.Getenv("MQTT_KEEP_ALIVE"); mqttKeepAlive != "" {
		v, _ := strconv.ParseInt(mqttKeepAlive, 10, 32)
		conf.MqttKeepAlive = int(v)
	}

	if mqttConnectTimeout := os.Getenv("MQTT_CONNECT_TIMEOUT"); mqttConnectTimeout != "" {
		v, _ := strconv.ParseInt(mqttConnectTimeout, 10, 32)
		conf.MqttConnectTimeout = int(v)
	}

	if mqttSessionsProvider := os.Getenv("MQTT_SESSIONS_PROVIDER"); mqttSessionsProvider != "" {
		conf.MqttSessionsProvider = mqttSessionsProvider
	}

	if mqttTopicsProvider := os.Getenv("MQTT_TOPICS_PROVIDER"); mqttTopicsProvider != "" {
		conf.MqttTopicsProvider = mqttTopicsProvider
	}

	if mqttUsername := os.Getenv("MQTT_USERNAME"); mqttUsername != "" {
		conf.MqttUsername = mqttUsername
	}

	if mqttPassword := os.Getenv("MQTT_PASSWORD"); mqttPassword != "" {
		conf.MqttPassword = mqttPassword
	}

	if mqttIp := os.Getenv("MQTT_IP"); mqttIp != "" {
		conf.MqttIp = mqttIp
	}

	if mqttPort := os.Getenv("MQTT_PORT"); mqttPort != "" {
		v, _ := strconv.ParseInt(mqttPort, 10, 32)
		conf.MqttPort = int(v)
	}

	if proxyPort := os.Getenv("PROXY_PORT"); proxyPort != "" {
		v, _ := strconv.ParseInt(proxyPort, 10, 32)
		conf.ProxyPort = int(v)
	}

	if serial := os.Getenv("SERIAL"); serial != "" {
		conf.Serial = strings.Split(serial, ",")
	}

}
