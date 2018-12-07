package smartbus

import (
	"time"
	"github.com/op/go-logging"
	"github.com/e154/smart-home-node/system/serial"
	"github.com/e154/smart-home-node/models/devices"
	"github.com/e154/smart-home-node/common"
	"encoding/json"
	"github.com/e154/smart-home-node/system/plugins/smartbus/driver"
)

var (
	log = logging.MustGetLogger("smartbus")
)

type Smartbus struct {
	params         *devices.DevSmartBusConfig

	command        []byte
	serialDev      *serial.Serial
	respFunc       func(data []byte)
	requestMessage *common.MessageRequest
}

func NewSmartbus(respFunc func(data []byte), requestMessage *common.MessageRequest) *Smartbus {

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

func (s *Smartbus) Open(dev string) (err error) {

	s.serialDev = &serial.Serial{
		Dev:         dev,
		Baud:        s.params.Baud,
		ReadTimeout: time.Duration(s.params.Timeout),
		StopBits:    s.params.StopBits,
	}

	if _, err = s.serialDev.Open(); err != nil {
		log.Warningf("%s - %s\r\n", dev, err.Error())
		return
	}

	return
}

func (s *Smartbus) Close() () {
	if s.serialDev != nil {
		s.serialDev.Close()
	}
	return
}

func (s *Smartbus) Exec(dev string) (resp *common.MessageResponse, err error) {

	resp = &common.MessageResponse{
		DeviceId:   s.requestMessage.DeviceId,
		DeviceType: s.requestMessage.DeviceType,
	}

	request := &devices.DevSmartBusRequest{}
	if err = json.Unmarshal(s.requestMessage.Command, request); err != nil {
		resp.Status = "error"
		return
	}

	// open
	if err = s.Open(dev); err != nil {
		resp.Status = "error"
		err = nil
		return
	}

	// exec command at port
	command := make([]byte, 0)
	command = append(command, byte(s.params.Device))
	command = append(command, request.Command...)

	modbus := &driver.Smartbus{Serial: s.serialDev}
	var result []byte
	if result, err = modbus.Send(command); err != nil {
		//errcode = "MODBUS_LINE_ERROR"
		log.Warningf("%s - %s\r\n", dev, err.Error())
		//TODO remove
		if err.Error() == "ILLEGAL_LRC" {
			err = nil
		} else {
			s.Close()
			return
		}
	}

	// bug in the devices need timeout, need fix!!!
	if s.params.Sleep != 0 {
		time.Sleep(time.Millisecond * time.Duration(s.params.Sleep))
	}
	s.Close()

	r := &devices.DevSmartBusResponse{
		Result: result,
	}
	if resp.Response, err = json.Marshal(r); err != nil {
		log.Error(err.Error())
	}

	return
}

func (s *Smartbus) Send(item interface{}) {
	data, _ := json.Marshal(item)
	s.respFunc(data)
}

func (s *Smartbus) DeviceId() int64 {
	return s.requestMessage.DeviceId
}
