package client

import (
	"sync"
	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/system/smartbus"
	"encoding/json"
	"github.com/e154/smart-home-node/models/devices"
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

func (t *Thread) Send(message *common.MessageRequest) (resp *common.MessageResponse, err error) {

	t.Lock()
	t.Busy = true
	defer func() {
		t.Unlock()
		t.Busy = false
	}()

	resp = &common.MessageResponse{
		DeviceId:   message.DeviceId,
		DeviceType: message.DeviceType,
	}

	switch message.DeviceType {
	// smartbus line
	case common.DevTypeSmartBus:
		params := &devices.DevSmartBusConfig{}
		json.Unmarshal(message.Properties, params)
		request := &devices.DevSmartBusRequest{}
		if err = json.Unmarshal(message.Command, request); err != nil {
			resp.Status = "error"
			return
		}
		bus := smartbus.NewSmartbus(message.DeviceId, params, t.Dev, request.Command)
		if _, err, _ = bus.Open(); err != nil {
			resp.Status = "error"
			err = nil
			return
		}
		var result []byte
		if result, err, _ = bus.Exec(); err != nil {
			resp.Status = "error"
			err = nil
		} else {
			resp.Status = "success"
		}
		bus.Close()

		r := &devices.DevSmartBusResponse{
			Result: result,
		}
		data, _ := json.Marshal(r)
		resp.Response = data

	default:
		log.Errorf("unknown device type %s", message.DeviceType)
	}

	return
}
