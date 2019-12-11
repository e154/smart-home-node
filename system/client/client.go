package client

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/system/cache"
	"github.com/e154/smart-home-node/system/config"
	"github.com/e154/smart-home-node/system/graceful_service"
	"github.com/e154/smart-home-node/system/mqtt"
	"github.com/e154/smart-home-node/system/mqtt_client"
	"github.com/e154/smart-home-node/system/plugins/command"
	"github.com/e154/smart-home-node/system/plugins/modbus"
	"github.com/e154/smart-home-node/system/plugins/smartbus"
	"github.com/e154/smart-home-node/system/serial"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
	"github.com/paulbellamy/ratecounter"
	"sync"
	"time"
)

const (
	threadTimeTick = 5 * time.Second
	pingTimeTick   = 1 * time.Second
)

var (
	log = logging.MustGetLogger("client")
)

type Client struct {
	Stat
	cfg                 *config.AppConfig
	mqtt                *mqtt.Mqtt
	mqttClient          *mqtt_client.Client
	updateThreadsTicker *time.Ticker
	updatePinkTicker    *time.Ticker
	status              common.ClientStatus
	cache               *cache.Cache
	startedAt           time.Time
	serialService       *serial.SerialService
	sync.Mutex
	pool Threads
}

func NewClient(cfg *config.AppConfig, graceful *graceful_service.GracefulService,
	mqtt *mqtt.Mqtt, serialService *serial.SerialService) *Client {

	memCache := &cache.Cache{
		Cachetime: 3600,
		Name:      "node",
	}
	client := &Client{
		cfg:                 cfg,
		updateThreadsTicker: time.NewTicker(threadTimeTick),
		updatePinkTicker:    time.NewTicker(pingTimeTick),
		pool:                make(Threads),
		cache:               memCache,
		status:              common.StatusEnabled,
		startedAt:           time.Now(),
		serialService:       serialService,
		Stat: Stat{
			rpsCounter: ratecounter.NewRateCounter(1 * time.Second),
			avgRequest: ratecounter.NewAvgRateCounter(60 * time.Second),
		},
		mqtt: mqtt,
	}

	graceful.Subscribe(client)

	go func() {
		for {
			select {
			case <-client.updatePinkTicker.C:
				client.ping()
			case <-client.updateThreadsTicker.C:
				client.UpdateThreads()
			}
		}
	}()

	return client
}

func (c *Client) Shutdown() {
	c.updateThreadsTicker.Stop()
	c.updatePinkTicker.Stop()
}

func (c *Client) Connect() {

	mqttClient, err := c.mqtt.NewClient(nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if err := mqttClient.Connect(); err != nil {
		log.Error(err.Error())
	}

	c.mqttClient = mqttClient

	_ = c.mqttClient.Subscribe(c.topic("req"), 0, c.onPublish)
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
	// smartbus plugin
	case common.DevTypeSmartBus:
		cmd := smartbus.NewSmartbus(c.ResponseFunc(cli), message)
		c.SendMessageToThread(cmd)
	// modbus rtu
	case common.DevTypeModBusRtu:
		cmd := modbus.NewModbusRtu(c.ResponseFunc(cli), message)
		c.SendMessageToThread(cmd)
	// modbus tcp
	case common.DevTypeModBusTcp:
		cmd := modbus.NewModbusTcp(c.ResponseFunc(cli), message)
		cmd.Exec()
	default:
		log.Warningf("unknown message device type: %s", message.DeviceType)
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
	cacheKey := c.cache.GetKey(fmt.Sprintf("%d_dev", item.DeviceId()))
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

	item.Send(resp)

	return
}

func (c *Client) UpdateThreads() {

	c.Lock()
	defer c.Unlock()

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

	if c.mqttClient != nil && (c.mqttClient.IsConnected()) {
		message := &common.ClientStatusModel{
			Status:    c.status,
			Thread:    activeThreads,
			Rps:       c.rpsCounter.Rate(),
			Min:       c.min,
			Max:       c.max,
			StartedAt: c.startedAt,
		}
		data, _ := json.Marshal(message)
		c.mqttClient.Publish(c.topic("ping"), data)
	}
}

func (c *Client) ResponseFunc(cli MQTT.Client) func(data []byte) {

	return func(data []byte) {
		// response
		if cli.IsConnected() {
			if token := cli.Publish(c.topic("resp"), 0x0, false, data); token.Wait() && token.Error() != nil {
				log.Error(token.Error().Error())
			}
		}
	}
}

func (c *Client) topic(r string) string {
	return fmt.Sprintf("/home/node/%s/%s", c.cfg.MqttClientId, r)
}