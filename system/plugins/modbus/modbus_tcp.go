package modbus

import (
	"encoding/binary"
	"encoding/json"
	"github.com/e154/smart-home-node/common"
	. "github.com/e154/smart-home-node/models/devices"
	"github.com/e154/smart-home-node/system/plugins/modbus/driver"
)

type ModbusTcp struct {
	params *DevModBusTcpConfig

	command        []byte
	respFunc       func(data []byte)
	requestMessage *common.MessageRequest
}

func NewModbusTcp(respFunc func(data []byte), requestMessage *common.MessageRequest) *ModbusTcp {

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
		DeviceId:   s.requestMessage.DeviceId,
		DeviceType: s.requestMessage.DeviceType,
		Status:     "success",
	}

	r := &DevModBusResponse{}

	request := &DevModBusRequest{}
	if err = json.Unmarshal(s.requestMessage.Command, request); err != nil {
		resp.Status = "error"
		return
	}

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

	//fmt.Println(string(q))

	s.respFunc(q)

	return
}

func (s *ModbusTcp) Send(item interface{}) {
	data, _ := json.Marshal(item)
	s.respFunc(data)
}

func (s *ModbusTcp) DeviceId() int64 {
	return s.requestMessage.DeviceId
}