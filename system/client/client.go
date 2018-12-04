package client

import (
	"fmt"
	"github.com/e154/smart-home-node/system/config"
	"github.com/e154/smart-home-node/system/graceful_service"
	"github.com/op/go-logging"
	"github.com/e154/smart-home-node/system/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"encoding/json"
	"time"
	"github.com/e154/smart-home-node/system/serial"
	"sync"
	"github.com/e154/smart-home-node/system/cache"
	"github.com/e154/smart-home-node/common"
	"github.com/paulbellamy/ratecounter"
	"github.com/e154/smart-home-node/system/command"
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
	client              *mqtt.Client
	updateThreadsTicker *time.Ticker
	updatePinkTicker    *time.Ticker
	status              common.ClientStatus
	cache               *cache.Cache
	startedAt           time.Time
	sync.Mutex
	pool Threads
}

func NewClient(cfg *config.AppConfig,
	graceful *graceful_service.GracefulService,
	qService *mqtt.Mqtt) *Client {

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
		Stat: Stat{
			rpsCounter: ratecounter.NewRateCounter(1 * time.Second),
			avgRequest: ratecounter.NewAvgRateCounter(60 * time.Second),
		},
	}
	topic := fmt.Sprintf("/home/%s", cfg.Topic)
	c, err := qService.NewClient(topic, 0x0, client.onPublish)
	if err != nil {
		log.Error(err.Error())
	}
	client.client = c

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
	c.client.Disconnect()
}

func (c *Client) Connect() {
	go c.client.Connect()
}

func (c *Client) onPublish(cli MQTT.Client, msg MQTT.Message) {

	c.rpsCounterIncr()

	message := &MessageReq{}
	if err := json.Unmarshal(msg.Payload(), message); err != nil {
		log.Error(err.Error())
		return
	}

	c.avgStart()
	switch message.DeviceType {
	case common.DevTypeCommand:
		command.NewCommand(cli, message.Command)
	case common.DevTypeSmartBus:
		c.SendMessageToThread(cli, message)
	default:
		log.Warningf("unknown message device type: %s", message.DeviceType)
	}

	c.avgEnd()

}

func (c *Client) SendMessageToThread(cli MQTT.Client, message *MessageReq) (err error) {

	//поиск в кэше
	cacheKey := c.cache.GetKey(fmt.Sprintf("%d_dev", message.DeviceId))
	var threadDev string
	if c.cache.IsExist(cacheKey) {
		threadDev = c.cache.Get(cacheKey).(string)
	}

	var resp *MessageResp
	if threadDev != "" {
		resp, err = c.pool[threadDev].Send(cli, message)
		return
	}

	for threadDev, thread := range c.pool {
		if resp, err = thread.Send(cli, message); err == nil {
			c.cache.Put("node", cacheKey, threadDev)
		}
	}

	// response
	if cli.IsConnected() {
		topic := fmt.Sprintf("/home/%s", c.cfg.Topic)
		data, _ := json.Marshal(resp)
		if token := cli.Publish(topic+"/resp", 0x0, false, data); token.Wait() && token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}

	return
}

func (c *Client) UpdateThreads() {

	c.Lock()
	defer c.Unlock()

	//log.Debug("update thread list")

	deviceList := serial.DeviceList()

	//remove threads
	for k := range c.pool {
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
		delete(c.pool, k)
	}

	//add threads
	for _, dev := range deviceList {
		if _, ok := c.pool[dev]; ok {
			continue
		}

		log.Debugf("Add thread to pool: %s", dev)
		c.pool[dev] = NewThread(dev)
	}

	if len(c.pool) == 0 {
		c.status = common.StatusBusy
	} else {
		c.status = common.StatusEnabled
	}
}

func (c *Client) ping() {
	if c.client != nil && (c.client.IsConnected()) {
		message := &ClientStatusModel{
			Status:    c.status,
			Thread:    len(c.pool),
			Rps:       c.rpsCounter.Rate(),
			Min:       c.min,
			Max:       c.max,
			StartedAt: c.startedAt,
		}
		data, _ := json.Marshal(message)
		c.client.Publish("/ping", data)
	}
}
