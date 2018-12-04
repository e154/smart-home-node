package client

import (
	"sync"
	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/system/smartbus"
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Thread struct {
	sync.Mutex
	Busy bool
	Dev  string
}

type Threads map[string]*Thread

func NewThread(dev string) (thread *Thread) {

	thread = &Thread{
		Dev:  dev,
		Busy: false,
	}

	return
}

func (t *Thread) Send(cli MQTT.Client, message *MessageReq) (resp *MessageResp, err error) {

	t.Lock()
	t.Busy = true
	defer func() {
		t.Unlock()
		t.Busy = false
	}()

	//debug.Println(message)

	resp = &MessageResp{
		DeviceId:   message.DeviceId,
		DeviceType: message.DeviceType,
	}

	switch message.DeviceType {
	// smartbus line
	case common.DevTypeSmartBus:
		params := &common.DevConfSmartBus{}
		json.Unmarshal(message.Properties, params)
		bus := smartbus.NewSmartbus(message.DeviceId, params, t.Dev, message.Command)
		if _, err, _ = bus.Open(); err != nil {
			resp.Status = "error"
			err = nil
			return
		}
		if resp.Response, err, _ = bus.Exec(); err != nil {
			bus.Close()
			resp.Status = "error"
			err = nil
			return
		}
		bus.Close()

		//fmt.Println(resp.Response)

	default:
		log.Errorf("unknown device type %s", message.DeviceType)
	}

	return
}
