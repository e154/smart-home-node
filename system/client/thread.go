package client

import (
	"sync"
	"github.com/e154/smart-home-node/common/debug"
	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/system/smartbus"
	"encoding/json"
	"fmt"
)

type Thread struct {
	sync.Mutex
	Busy bool
	Dev  string
	d int
}

type Threads map[string]*Thread

func NewThread(dev string) (thread *Thread) {

	thread = &Thread{
		Dev:  dev,
		Busy: false,
	}

	return
}

func (t *Thread) Send(message *MessageReq) (resp *MessageResp, err error) {

	t.d++

	t.Lock()
	t.Busy = true
	defer func() {
		t.Unlock()
		t.Busy = false
	}()

	debug.Println(message)

	resp = &MessageResp{
		DeviceId:   message.DeviceId,
		DeviceType: message.DeviceType,
	}

	switch message.DeviceType {
	// smartbus line
	case common.DevTypeSmartBus:
		params := &common.DevConfSmartBus{}
		json.Unmarshal(message.Properties, params)
		params.Device = t.d
		bus := smartbus.NewSmartbus(message.DeviceId, params, t.Dev, message.Command)
		if _, err, _ = bus.Open(); err != nil {
			return
		}
		var res string
		if res, err, _ = bus.Exec(); err != nil {
			bus.Close()
			return
		}
		bus.Close()

		fmt.Println(res)
	default:
		log.Errorf("unknown device type %s", message.DeviceType)
	}

	return
}
