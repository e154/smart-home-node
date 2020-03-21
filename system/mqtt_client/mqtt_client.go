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

package mqtt_client

import (
	"github.com/e154/smart-home-node/common"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var (
	log = common.MustGetLogger("mqtt_client")
)

type Client struct {
	cfg            *Config
	client         MQTT.Client
	subscribesLock *sync.Mutex
	subscribes     map[string]Subscribe
	inConnectProc  *atomic.Bool
	disconnect     *atomic.Bool
}

func NewClient(cfg *Config) (client *Client, err error) {

	log.Infof("new queue client(%s) uri(%s)", cfg.ClientID, cfg.Broker)

	client = &Client{
		cfg:            cfg,
		subscribesLock: &sync.Mutex{},
		subscribes:     make(map[string]Subscribe),
		disconnect:     atomic.NewBool(false),
		inConnectProc:  atomic.NewBool(false),
	}

	opts := MQTT.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID).
		SetKeepAlive(time.Duration(cfg.KeepAlive) * time.Second).
		SetPingTimeout(time.Duration(cfg.PingTimeout) * time.Second).
		SetConnectTimeout(time.Duration(cfg.ConnectTimeout) * time.Second).
		SetCleanSession(cfg.CleanSession).
		SetUsername(cfg.Username).
		SetPassword(cfg.Password).
		SetOnConnectHandler(client.onConnect).
		SetConnectionLostHandler(client.onConnectionLostHandler)

	client.client = MQTT.NewClient(opts)

	return
}

func (c *Client) Connect() {

	if c.inConnectProc.Load() {
		return
	}

	c.inConnectProc.Store(true)


	log.Infof("Connect to server %s", c.cfg.Broker)

	go func() {
		defer func() {
			c.inConnectProc.Store(false)
		}()

	LOOP1:
		token := c.client.Connect()

	LOOP2:
		token.Wait()

		if err := token.Error(); err != nil {
			time.Sleep(time.Second)
			goto LOOP1
		}

		if c.disconnect.Load() {
			return
		}

		time.Sleep(time.Second)
		goto LOOP2
	}()
}

func (c *Client) Disconnect() {
	if c.client == nil {
		return
	}

	c.disconnect.Store(true)
	c.safeUnsubscribeAll()
	c.subscribesLock.Lock()
	c.subscribes = make(map[string]Subscribe)
	c.subscribesLock.Unlock()
	c.client.Disconnect(250)
	c.client = nil
}

func (c *Client) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) (err error) {

	c.subscribesLock.Lock()
	defer c.subscribesLock.Unlock()

	if _, ok := c.subscribes[topic]; !ok {
		c.subscribes[topic] = Subscribe{
			Qos:      qos,
			Callback: callback,
		}
	}

	if token := c.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
		return token.Error()
	}
	return
}

func (c *Client) Unsubscribe(topic string) (err error) {

	if token := c.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
		return token.Error()
	}
	return
}

func (c *Client) safeUnsubscribeAll() {

	c.subscribesLock.Lock()
	defer c.subscribesLock.Unlock()

	for topic, _ := range c.subscribes {
		if token := c.client.Unsubscribe(topic); token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
	c.subscribes = make(map[string]Subscribe)
}

func (c *Client) Publish(topic string, payload interface{}) (err error) {
	if c.client != nil && (c.client.IsConnected()) {
		c.client.Publish(topic, c.cfg.Qos, false, payload)
	}
	return
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnectionOpen()
}

func (c *Client) onConnectionLostHandler(client MQTT.Client, e error) {

	log.Debug("connection lost...")

	c.subscribesLock.Lock()
	defer c.subscribesLock.Unlock()

	for topic, _ := range c.subscribes {
		if token := c.client.Unsubscribe(topic); token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
}

func (c *Client) onConnect(client MQTT.Client) {

	log.Debug("connected...")

	c.subscribesLock.Lock()
	defer c.subscribesLock.Unlock()

	for topic, subscribe := range c.subscribes {
		if token := c.client.Subscribe(topic, subscribe.Qos, subscribe.Callback); token.Wait() && token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
}
