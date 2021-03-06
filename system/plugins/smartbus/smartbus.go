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

package smartbus

import (
	"encoding/json"
	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/models/devices"
	"github.com/e154/smart-home-node/system/plugins/smartbus/driver"
	"time"
)

var (
	log = common.MustGetLogger("smartbus")
)

type Smartbus struct {
	params *devices.DevSmartBusConfig

	command        []byte
	respFunc       func(deviceId int64, data []byte)
	requestMessage *common.MessageRequest
}

func NewSmartbus(respFunc func(deviceId int64, data []byte), requestMessage *common.MessageRequest) *Smartbus {

	params := &devices.DevSmartBusConfig{}
	if err := json.Unmarshal(requestMessage.Properties, params); err != nil {
		log.Error(err.Error())
	}

	return &Smartbus{
		params:         params,
		command:        requestMessage.Command,
		respFunc:       respFunc,
		requestMessage: requestMessage,
	}
}

func (s *Smartbus) Exec(t common.Thread) (resp *common.MessageResponse, err error) {

	resp = &common.MessageResponse{
		DeviceId:   s.requestMessage.DeviceId,
		DeviceType: s.requestMessage.DeviceType,
		Status:     "success",
	}

	r := &devices.DevSmartBusResponse{}

	request := &devices.DevSmartBusRequest{}
	if err = json.Unmarshal(s.requestMessage.Command, request); err != nil {
		resp.Status = "error"
		return
	}

	// open
	if err = t.SetParams(s.params.Baud, s.params.Timeout, s.params.StopBits); err != nil {
		resp.Status = "error"
		err = nil
		return
	}

	// exec command at port
	command := make([]byte, 0)
	command = append(command, byte(s.params.Device))
	command = append(command, request.Command...)

	modbus := &driver.Smartbus{Serial: t.GetSerial()}
	if r.Result, err = modbus.Send(command); err != nil {
		t.SetErr()

		//errcode = "MODBUS_LINE_ERROR"
		log.Warnf("%s - %s\r\n", t.Device(), err.Error())
		//TODO remove
		if err.Error() == "ILLEGAL_LRC" {
			err = nil
		} else {
			r.Error = err.Error()
			resp.Status = "error"
		}
	}

	// bug in the devices need timeout, need fix!!!
	if s.params.Sleep != 0 {
		time.Sleep(time.Millisecond * time.Duration(s.params.Sleep))
	}

	if resp.Response, err = json.Marshal(r); err != nil {
		log.Error(err.Error())
	}

	return
}

func (s *Smartbus) Send(deviceId int64, item interface{}) {
	data, _ := json.Marshal(item)
	s.respFunc(deviceId, data)
}

func (s *Smartbus) DeviceId() int64 {
	return s.requestMessage.DeviceId
}
