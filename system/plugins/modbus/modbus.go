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

	var firstTime bool

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
		//log.Error(err.Error())
		r.Error = err.Error()
		if firstTime {
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

	time.Sleep(time.Millisecond * 10)

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
