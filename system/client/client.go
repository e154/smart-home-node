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
)

const (
	TimeTick = 30 * time.Second
)

var (
	log = logging.MustGetLogger("client")
)

type Client struct {
	cfg                 *config.AppConfig
	client              *mqtt.Client
	updateThreadsTicker *time.Ticker
	sync.Mutex
	pool  Threads
	cache *cache.Cache
}

func NewClient(cfg *config.AppConfig,
	graceful *graceful_service.GracefulService,
	qService *mqtt.Mqtt) *Client {

	cache := &cache.Cache{
		Cachetime: 3600,
		Name:      "node",
	}
	client := &Client{
		cfg:                 cfg,
		updateThreadsTicker: time.NewTicker(TimeTick),
		pool:                make(Threads),
		cache:               cache,
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
			case <-client.updateThreadsTicker.C:
				client.UpdateThreads()
			}
		}
	}()

	return client
}

func (c *Client) Shutdown() {
	c.updateThreadsTicker.Stop()
	c.client.Disconnect()
}

func (c *Client) Connect() {
	go c.client.Connect()
}

func (c *Client) onPublish(cli MQTT.Client, msg MQTT.Message) {

	message := &MessageReq{}
	if err := json.Unmarshal(msg.Payload(), message); err != nil {
		log.Error(err.Error())
	}

	resp, err := c.SendMessageToThread(message)
	if err != nil {
		log.Error(err.Error())
	}

	// response
	if cli.IsConnected() {
		topic := fmt.Sprintf("/home/%s", c.cfg.Topic)
		data, _ := json.Marshal(resp)
		if token := cli.Publish(topic+"/resp", 0x0, false, data); token.Wait() && token.Error() != nil {
			log.Error(token.Error().Error())
		}
	}
}

func (c *Client) SendMessageToThread(message *MessageReq) (resp *MessageResp, err error) {

	//поиск в кэше
	cacheKey := c.cache.GetKey(fmt.Sprintf("%d_dev", message.DeviceId))
	var dev string
	if c.cache.IsExist(cacheKey) {
		dev = c.cache.Get(cacheKey).(string)
	}

	if dev != "" {
		resp, err = c.pool[dev].Send(message)
		return
	}

	for _, thread := range c.pool {
		if resp, err = thread.Send(message); err == nil {
			c.cache.Put("node", cacheKey, message.DeviceId)
		}
	}

	return
}

func (c *Client) UpdateThreads() {

	c.Lock()
	defer c.Unlock()

	log.Debug("update thread list")

	deviceList := serial.DeviceList()

	//remove threads
	for k, _ := range c.pool {
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
}
