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

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/common/logger"
	"github.com/e154/smart-home-node/system/cache"
	"github.com/e154/smart-home-node/system/config"
	"github.com/e154/smart-home-node/system/mqtt"
	"github.com/e154/smart-home-node/system/mqtt_client"
	"github.com/e154/smart-home-node/system/plugins/command"
	"github.com/e154/smart-home-node/system/plugins/modbus"
	"github.com/e154/smart-home-node/system/serial"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/fx"
)

const (
	threadTimeTick = 5 * time.Second
	pingTimeTick   = 1 * time.Second
)

var (
	log = logger.MustGetLogger("client")
)

type Client struct {
	Stat
	cfg                 *config.AppConfig
	mqtt                *mqtt.Mqtt
	updateThreadsTicker *time.Ticker
	updatePinkTicker    *time.Ticker
	status              common.ClientStatus
	cache               *cache.Cache
	serialService       *serial.SerialService
	poolLocker          sync.Mutex
	pool                Threads
	mqttClientLocker    sync.Mutex
	mqttClient          *mqtt_client.Client
}

func NewClient(lc fx.Lifecycle,
	cfg *config.AppConfig,
	mqtt *mqtt.Mqtt,
	serialService *serial.SerialService) *Client {

	memCache := &cache.Cache{
		Cachetime: 3600,
		Name:      "node",
	}
	client := &Client{
		cfg:              cfg,
		cache:            memCache,
		status:           common.StatusEnabled,
		serialService:    serialService,
		Stat:             NewStat(),
		mqtt:             mqtt,
		poolLocker:       sync.Mutex{},
		pool:             make(Threads),
		mqttClientLocker: sync.Mutex{},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return client.Start()
		},
		OnStop: func(ctx context.Context) (err error) {
			return client.Shutdown()
		},
	})

	return client
}

func (c *Client) Start() error {

	c.updateThreadsTicker = time.NewTicker(threadTimeTick)
	c.updatePinkTicker = time.NewTicker(pingTimeTick)

	go func() {
		for {
			select {
			case <-c.updatePinkTicker.C:
				c.ping()
			case <-c.updateThreadsTicker.C:
				c.UpdateThreads()
			}
		}
	}()

	c.Connect()

	return nil
}

func (c *Client) Shutdown() error {
	if c.updateThreadsTicker != nil {
		c.updateThreadsTicker.Stop()
	}
	if c.updatePinkTicker != nil {
		c.updatePinkTicker.Stop()
	}
	return nil
}

func (c *Client) Connect() {

	c.mqttClientLocker.Lock()
	defer c.mqttClientLocker.Unlock()

	var err error
	if c.mqttClient, err = c.mqtt.NewClient(nil); err != nil {
		log.Error(err.Error())
	}

	_ = c.mqttClient.Subscribe(c.topic("req/#"), 0, c.onPublish)

	c.mqttClient.Connect()
}

func (c *Client) onPublish(cli MQTT.Client, msg MQTT.Message) {

	c.rpsCounterIncr()

	message := &common.MessageRequest{}
	if err := json.Unmarshal(msg.Payload(), message); err != nil {
		log.Error(err.Error())
		return
	}

	startTime := c.avgStart()
	switch message.DeviceType {
	// command plugin
	case common.DevTypeCommand:
		cmd := command.NewCommand(c.ResponseFunc(cli), message)
		cmd.Exec()
	// modbus rtu
	case common.DevTypeModBusRtu:
		cmd := modbus.NewModbusRtu(c.ResponseFunc(cli), message)
		c.SendMessageToThread(cmd)
	// modbus tcp
	case common.DevTypeModBusTcp:
		cmd := modbus.NewModbusTcp(c.ResponseFunc(cli), message)
		cmd.Exec()
	default:
		log.Warnf("unknown message device type: %s", message.DeviceType)
	}

	c.avgEnd(startTime)

}

func (c *Client) SendMessageToThread(item common.ThreadCaller) (err error) {

	var activeThreads int
	for _, thread := range c.pool {
		if thread.Active {
			activeThreads++
		}
	}

	if activeThreads == 0 {
		return
	}

LOOP:
	//поиск в кэше
	cacheKey := c.cache.GetKey(fmt.Sprintf("%s_entityId", item.EntityId()))
	var threadDev string
	if c.cache.IsExist(cacheKey) {
		threadDev = c.cache.Get(cacheKey).(string)
	}

	var resp *common.MessageResponse
	if threadDev != "" {
		if c.pool[threadDev].Active {
			if resp, err = c.pool[threadDev].Exec(item); err != nil {
				c.cache.Delete(cacheKey)
			}
		} else {
			c.cache.Delete(cacheKey)
			goto LOOP
		}
	} else {
		for threadDev, thread := range c.pool {
			if !thread.Active {
				continue
			}
			if resp, err = thread.Exec(item); err == nil {
				c.cache.Put("node", cacheKey, threadDev)
				break
			}
		}
	}

	if err != nil {
		return
	}

	item.Send(item.EntityId(), resp)

	return
}

func (c *Client) UpdateThreads() {

	c.poolLocker.Lock()
	defer c.poolLocker.Unlock()

	//log.Debug("update thread list")

	deviceList := c.serialService.DeviceList()

	//remove threads
	for k, thread := range c.pool {
		if !thread.Active {
			continue
		}

		var exist bool
		for _, dev := range deviceList {
			if dev == k {
				exist = true
			}
		}
		if exist {
			continue
		}

		log.Debugf("Remove thread from pool: %s", k)
		thread.Disable()
		c.cache.ClearGroup("node")
		//delete(c.pool, k)
	}

	//add threads
	for _, dev := range deviceList {
		if thread, ok := c.pool[dev]; ok {
			if !thread.Active {
				log.Debugf("Add thread to pool: %s", dev)
				thread.Enable()
			}
			continue
		}

		log.Debugf("Add thread to pool: %s", dev)
		c.pool[dev] = NewThread(dev)
		c.cache.ClearGroup("node")
	}

	// check active threads
	var activeThreads int
	for _, thread := range c.pool {
		if thread.Active {
			activeThreads++
		}
	}

	c.status = common.StatusEnabled
}

func (c *Client) ping() {
	var activeThreads int
	for _, thread := range c.pool {
		if thread.Active {
			activeThreads++
		}
	}

	c.mqttClientLocker.Lock()
	defer c.mqttClientLocker.Unlock()

	if c.mqttClient != nil && (c.mqttClient.IsConnected()) {
		snapshot := c.GetStat()
		message := common.ClientStatusModel{
			Status:    c.status,
			Thread:    activeThreads,
			Rps:       snapshot.Rps,
			Min:       snapshot.Min,
			Max:       snapshot.Max,
			Latency:   snapshot.Latency,
			StartedAt: snapshot.StartedAt,
		}
		data, _ := json.Marshal(message)
		c.mqttClient.Publish(c.topic("ping"), data)
	}
}

func (c *Client) ResponseFunc(cli MQTT.Client) func(entityId string, data []byte) {

	return func(entityId string, data []byte) {
		// response
		if cli.IsConnected() {
			if token := cli.Publish(c.topic(fmt.Sprintf("resp/%s", entityId)), 0x0, false, data); token.Wait() && token.Error() != nil {
				log.Error(token.Error().Error())
			}
		}
	}
}

func (c *Client) topic(r string) string {
	return fmt.Sprintf("home/node/%s/%s", c.cfg.MqttClientId, r)
}
