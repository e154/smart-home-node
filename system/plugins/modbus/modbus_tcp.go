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

	"github.com/e154/smart-home-node/common"
	. "github.com/e154/smart-home-node/models/devices"
	modbus "github.com/e154/smart-home-node/system/plugins/modbus/driver"
)

type ModbusTcp struct {
	params *DevModBusTcpConfig

	command        []byte
	respFunc       func(entityId string, data []byte)
	requestMessage *common.MessageRequest
}

func NewModbusTcp(respFunc func(entityId string, data []byte), requestMessage *common.MessageRequest) *ModbusTcp {

	params := &DevModBusTcpConfig{}
	if err := json.Unmarshal(requestMessage.Properties, params); err != nil {
		log.Error(err.Error())
	}

	return &ModbusTcp{
		params:         params,
		command:        requestMessage.Command,
		respFunc:       respFunc,
		requestMessage: requestMessage,
	}
}

func (s *ModbusTcp) Exec() (resp *common.MessageResponse, err error) {

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

	log.Debugf("command func(%s) address(%d), count(%d), command(%v)", request.Function, request.Address, request.Count, request.Command)

	// set value
	value := make([]byte, 0)
	v := make([]byte, 2)
	for _, item := range request.Command {
		binary.BigEndian.PutUint16(v, item)
		value = append(value, v...)
	}

	cli := modbus.NewClient(modbus.NewTCPClientHandler(s.params.AddressPort))
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

	resp.Response, _ = json.Marshal(r)
	q, _ := json.Marshal(resp)

	//log.Debugf("result: %v", r.Result)

	s.respFunc(s.requestMessage.EntityId, q)

	return
}

func (s *ModbusTcp) Send(entityId string, item interface{}) {
	data, _ := json.Marshal(item)
	s.respFunc(entityId, data)
}

func (s *ModbusTcp) EntityId() string {
	return s.requestMessage.EntityId
}
