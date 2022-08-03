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

package modbus

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/common/logger"
	. "github.com/e154/smart-home-node/models/devices"
	modbus "github.com/e154/smart-home-node/system/plugins/modbus/driver"
)

var (
	log = logger.MustGetLogger("modbus")
)

type ModbusRtu struct {
	params *DevModBusRtuConfig

	command        []byte
	respFunc       func(entityId string, data []byte)
	requestMessage *common.MessageRequest
}

func NewModbusRtu(respFunc func(entityId string, data []byte), requestMessage *common.MessageRequest) *ModbusRtu {

	params := &DevModBusRtuConfig{}
	if err := json.Unmarshal(requestMessage.Properties, params); err != nil {
		log.Error(err.Error())
	}

	return &ModbusRtu{
		params:         params,
		command:        requestMessage.Command,
		respFunc:       respFunc,
		requestMessage: requestMessage,
	}
}

func (s *ModbusRtu) Exec(t common.Thread) (resp *common.MessageResponse, err error) {

	//startTime := time.Now()
	//fmt.Println("exec <-----")
	//defer func() {
	//	total := time.Since(startTime).Seconds()
	//	fmt.Println("exit ----->", total)
	//}()

	var firstTime bool

	resp = &common.MessageResponse{
		EntityId:   s.requestMessage.EntityId,
		DeviceType: s.requestMessage.DeviceType,
		Status:     "success",
	}

	r := &DevModBusResponse{}

	request := &DevModBusRequest{}
	if err = json.Unmarshal(s.requestMessage.Command, request); err != nil {
		resp.Status = "error"
		return
	}

	//debug.Println(s.params)
	//fmt.Println("device ", t.Device())

	con := t.GetCon()
	var handler *modbus.RTUClientHandler

LOOP:
	if con == nil {
		firstTime = true
		if handler, err = s.Connect(t.Device()); err != nil {
			resp.Status = "error"
			return
		}

		t.SetCon(handler)

	} else {
		switch v := con.(type) {
		case *modbus.RTUClientHandler:
			handler = v
		default:
			log.Errorf("unknown con type %v", v)
			con = nil
			goto LOOP
		}
		s.Check(handler)
	}

	// set value
	value := make([]byte, 0)
	v := make([]byte, 2)
	for _, item := range request.Command {
		binary.BigEndian.PutUint16(v, item)
		value = append(value, v...)
	}

	cli := modbus.NewClient(handler)
	var results []byte
	switch request.Function {
	case ReadInputRegisters:
		results, err = cli.ReadInputRegisters(request.Address, request.Count)
	case ReadHoldingRegisters:
		results, err = cli.ReadHoldingRegisters(request.Address, request.Count)
	case WriteSingleRegister:
		// count as value
		results, err = cli.WriteSingleRegister(request.Address, request.Count)
	case WriteMultipleRegisters:
		results, err = cli.WriteMultipleRegisters(request.Address, request.Count, value)
	case ReadCoils:
		results, err = cli.ReadCoils(request.Address, request.Count)
	case ReadDiscreteInputs:
		results, err = cli.ReadDiscreteInputs(request.Address, request.Count)
	case WriteSingleCoil:
		// count as value
		results, err = cli.WriteSingleCoil(request.Address, request.Count)
	case WriteMultipleCoils:
		results, err = cli.WriteMultipleCoils(request.Address, request.Count, value)
	default:
		log.Errorf("unknown function %s", request.Function)
	}

	if err != nil {
		resp.Status = "error"
		log.Error(err.Error())
		r.Error = err.Error()
		if firstTime {
			//fmt.Println("clear handler", t.Device())
			handler.Close()
			t.SetCon(nil)
		}
	}

	k := 0
	for j := 0; j < len(results); j++ {
		if k > 0 {
			k = 0
			qw := binary.BigEndian.Uint16([]byte{results[j-1], results[j]})
			r.Result = append(r.Result, qw)
		} else {
			k++
		}
	}

	//fmt.Println(r.Result)
	//fmt.Println("---")

	resp.Response, _ = json.Marshal(r)

	return
}

func (s *ModbusRtu) Send(entityId string, item interface{}) {
	data, _ := json.Marshal(item)
	s.respFunc(entityId, data)
}

func (s *ModbusRtu) EntityId() string {
	return s.requestMessage.EntityId
}

func (s *ModbusRtu) Connect(device string) (handler *modbus.RTUClientHandler, err error) {

	handler = modbus.NewRTUClientHandler(device)
	handler.BaudRate = s.params.Baud
	handler.DataBits = s.params.DataBits
	handler.Parity = s.parity(s.params.Parity)
	handler.StopBits = s.params.StopBits
	handler.SlaveId = byte(s.params.SlaveId)
	handler.Timeout = time.Duration(s.params.Timeout) * time.Millisecond
	//handler.Logger = l12.New(os.Stdout, "test: ", l12.LstdFlags)
	//handler.IdleTimeout = 100 * time.Millisecond

	if err = handler.Connect(); err != nil {
		log.Error(err.Error())
		return
	}
	//defer func() {
	//	fmt.Println("close handler")
	//	handler.Close()
	//}()

	time.Sleep(time.Millisecond * 100)

	return
}

func (s *ModbusRtu) Check(handler *modbus.RTUClientHandler) {

	var restart bool
	if handler.BaudRate != s.params.Baud {
		restart = true
	}
	if handler.DataBits != s.params.DataBits {
		restart = true
	}
	if handler.Parity != s.parity(s.params.Parity) {
		restart = true
	}
	if handler.StopBits != s.params.StopBits {
		restart = true
	}

	if restart {
		handler.Close()
		time.Sleep(100 * time.Millisecond)
		handler, _ = s.Connect(handler.Address)
	}
}

func (s *ModbusRtu) parity(p string) (parity string) {
	switch p {
	case "odd":
		parity = "O"
	case "even":
		parity = "E"
	default:
		parity = "N"
	}
	return
}
