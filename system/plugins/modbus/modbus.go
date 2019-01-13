package modbus

import (
	"github.com/e154/smart-home-node/common"
	"encoding/json"
	. "github.com/e154/smart-home-node/models/devices"
	"github.com/op/go-logging"
	"github.com/goburrow/modbus"
	"time"
	"encoding/binary"
)

var (
	log = logging.MustGetLogger("modbus")
)

type Modbus struct {
	params *DevModBusConfig

	command        []byte
	respFunc       func(data []byte)
	requestMessage *common.MessageRequest
}

func NewModbus(respFunc func(data []byte), requestMessage *common.MessageRequest) *Modbus {

	params := &DevModBusConfig{}
	if err := json.Unmarshal(requestMessage.Properties, params); err != nil {
		log.Error(err.Error())
	}

	return &Modbus{
		params:         params,
		command:        requestMessage.Command,
		respFunc:       respFunc,
		requestMessage: requestMessage,
	}
}

func (s *Modbus) Exec(t common.Thread) (resp *common.MessageResponse, err error) {

	//startTime := time.Now()
	//fmt.Println("exec <-----")
	//defer func() {
	//	total := time.Since(startTime).Seconds()
	//	fmt.Println("exit ----->", total)
	//}()

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

	//debug.Println(s.params)
	//debug.Println(t.Device())
	con := t.GetCon()
LOOP:
	var handler *modbus.RTUClientHandler
	if con == nil {
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
			return
		}
		s.Check(handler)
	}

	time.Sleep(time.Millisecond * 10)

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
		results, err = cli.WriteMultipleRegisters(request.Address, request.Count, request.Command)
	case ReadCoils:
		results, err = cli.ReadCoils(request.Address, request.Count)
	case ReadDiscreteInputs:
		results, err = cli.ReadDiscreteInputs(request.Address, request.Count)
	case WriteSingleCoil:
		// count as value
		results, err = cli.WriteSingleCoil(request.Address, request.Count)
	case WriteMultipleCoils:
		results, err = cli.WriteMultipleCoils(request.Address, request.Count, request.Command)
	default:
		log.Errorf("unknown function %s", request.Function)
	}

	if err != nil {
		resp.Status = "error"
		log.Error(err.Error())
		return
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

	if resp.Response, err = json.Marshal(r); err != nil {
		log.Error(err.Error())
	}

	return
}

func (s *Modbus) Send(item interface{}) {
	data, _ := json.Marshal(item)
	s.respFunc(data)
}

func (s *Modbus) DeviceId() int64 {
	return s.requestMessage.DeviceId
}

func (s *Modbus) Connect(device string) (handler *modbus.RTUClientHandler, err error) {

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

	return
}

func (s *Modbus) Check(handler *modbus.RTUClientHandler) {

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

func (s *Modbus) parity(p string) (parity string) {
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
