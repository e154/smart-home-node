package smartbus

import (
	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/system/serial"
	"time"
	"github.com/e154/smart-home-node/system/smartbus/driver"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("smartbus")
)

type Smartbus struct {
	params    *common.DevConfSmartBus
	deviceId  int64
	dev       string
	command   []byte
	serialDev *serial.Serial
}

func NewSmartbus(deviceId int64, params *common.DevConfSmartBus, dev string, command []byte) *Smartbus {
	return &Smartbus{
		deviceId: deviceId,
		params:   params,
		dev:      dev,
		command:  command,
	}
}

func (s *Smartbus) Open() (result string, err error, errcode string) {

	s.serialDev = &serial.Serial{
		Dev:         s.dev,
		Baud:        s.params.Baud,
		ReadTimeout: time.Duration(s.params.Timeout),
		StopBits:    s.params.StopBits,
	}

	if _, err = s.serialDev.Open(); err != nil {
		errcode = "SERIAL_PORT_ERROR"
		log.Warningf("%s - %s\r\n", s.dev, err.Error())
		return
	}

	return
}

func (s *Smartbus) Close() (result string, err error, errcode string) {
	if s.serialDev != nil {
		s.serialDev.Close()
	}
	return
}

func (s *Smartbus) Exec() (result []byte, err error, errcode string) {

	command := make([]byte, 0)

	command = append(command, byte(s.params.Device))
	command = append(command, s.command...)

	modbus := &driver.Smartbus{Serial: s.serialDev}
	if result, err = modbus.Send(command); err != nil {
		errcode = "MODBUS_LINE_ERROR"
		log.Warningf("%s - %s\r\n", s.dev, err.Error())
		//TODO remove
		if err.Error() == "ILLEGAL_LRC" {
			err = nil
		} else {
			return
		}
	}

	// bug in the devices need timeout, need fix!!!
	if s.params.Sleep != 0 {
		time.Sleep(time.Millisecond * time.Duration(s.params.Sleep))
	}

	return
}
