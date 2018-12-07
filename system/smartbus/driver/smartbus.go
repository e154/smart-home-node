package driver

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/e154/smart-home-node/system/serial"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("smartbus/driver")
)

const (
	IllegalFunction uint8 = iota + 1
	IllegalDataAddress
	IllegalDataValue
	SlaveDeviceFailure
	acknowledge
	SlaveDeviceBusy
	NegativeAcknowledge
	MemoryParityError
	IllegalLrc
)

const (
	StateRxIdle uint8 = iota
	StateRxRcv
	StateRxWaitEof
)

const (
	ByteHighNibble uint8 = iota
	ByteLowNibble
)

type Smartbus struct {
	Serial       *serial.Serial
	rcvState     uint8
	rcvBytePos   uint8
	rcvBufferPos uint8
	rcvBuf       []byte
	trcBuff      *bytes.Buffer
}

// 1 - address		u08
// 1 - function		u08
// 1..n - data		u08 x N
func (m *Smartbus) Send(data []byte) (result []byte, err error) {

	var b byte
	var ok bool
	lrc := LRC(data)
	data = append(data, lrc)

	if err = m.asciiTransmit(data); err != nil {
		return
	}

	m.rcvState = StateRxIdle
	reader := bufio.NewReader(m.Serial.Port)
	for {
		if b, err = reader.ReadByte(); err != nil {
			break
		}

		if ok, err = m.asciiReceiveFSM(b); err != nil || ok {
			break
		}
	}

	if len(m.rcvBuf) > 2 {
		result = m.rcvBuf[2 : len(m.rcvBuf)-1]
	}

	log.Debugf("receive <- %X, len %d\r\n", result, len(m.rcvBuf)) //TODO remove

	return
}

func (m *Smartbus) asciiReceiveFSM(b byte) (bool, error) {

	//log.Print(string(b))
	switch m.rcvState {
	case StateRxRcv:
		if b == ':' {
			m.rcvBuf = []byte{}
			m.rcvBytePos = ByteHighNibble;
		} else if b == '\r' {
			m.rcvState = StateRxWaitEof;
		} else {
			b = char2bin(b)
			switch m.rcvBytePos {
			case ByteHighNibble:
				m.rcvBuf = append(m.rcvBuf, b<<4)
				m.rcvBytePos = ByteLowNibble
			case ByteLowNibble:
				m.rcvBuf[len(m.rcvBuf)-1] |= b
				m.rcvBytePos = ByteHighNibble
			}
		}
	case StateRxWaitEof:
		if b == '\n' {
			m.rcvState = StateRxIdle
			log.Debugf("receive <- %X, len: %d\r\n", m.rcvBuf, len(m.rcvBuf)) //TODO remove
			return true, checkError(m.rcvBuf)

		} else if b == ':' {
			m.rcvBuf = []byte{}
			m.rcvBytePos = ByteHighNibble;
			m.rcvState = StateRxRcv;
		} else {
			m.rcvState = StateRxIdle;
		}
	case StateRxIdle:
		if b == ':' {
			m.rcvBuf = []byte{}
			m.rcvBytePos = ByteHighNibble;
			m.rcvState = StateRxRcv;
		}
	}

	return false, nil
}

// 1 - address		u08
// 1 - function		u08
// 1..n - data		u08 x N
// 1 - lrc			u08
// 1 - \r			u08
// 1 - \n			u08
func (m *Smartbus) asciiTransmit(data []byte) (err error) {

	if m.trcBuff != nil {
		m.trcBuff.Reset()
	}

	m.trcBuff = &bytes.Buffer{}
	if err = m.trcBuff.WriteByte(':'); err != nil {
		return
	}

	for _, d := range data {
		m.trcBuff.WriteByte(bin2char(HI(d)))
		m.trcBuff.WriteByte(bin2char(LOW(d)))
	}

	m.trcBuff.Write([]byte{'\r', '\n'})

	if m.Serial == nil {
		err = errors.New("serial pointer is nil")
		return
	}

	//log.Debugf("send -> %X, %v\r\n", m.trcBuff.Bytes(), data) //TODO comment

	_, err = m.Serial.Port.Write(m.trcBuff.Bytes())
	if err != nil {
		return
	}

	return
}

func checkError(buf []byte) error {

	var errCode uint8

	// check lrc
	if LRC(buf[0:len(buf)-1]) != buf[len(buf)-1] {
		errCode = IllegalLrc
	}

	// check error bite
	if buf[1]&(1<<7) != 0 {
		errCode = buf[2]
	}

	// convert error code
	var err error
	switch errCode {
	case IllegalFunction:
		err = errors.New("ILLEGAL_FUNCTION")
	case IllegalDataAddress:
		err = errors.New("ILLEGAL_DATA_ADDRESS")
	case IllegalDataValue:
		err = errors.New("ILLEGAL_DATA_VALUE")
	case SlaveDeviceFailure:
		err = errors.New("SLAVE_DEVICE_FAILURE")
	case acknowledge:
		err = errors.New("ACKNOWLEDGE")
	case SlaveDeviceBusy:
		err = errors.New("SLAVE_DEVICE_BUSY")
	case NegativeAcknowledge:
		err = errors.New("NEGATIVE_ACKNOWLEDGE")
	case MemoryParityError:
		err = errors.New("MEMORY_PARITY_ERROR")
	case IllegalLrc:
		err = errors.New("ILLEGAL_LRC")
	}

	return err
}

func bin2char(b byte) byte {

	if b <= 0x09 {
		return byte('0' + b)
	} else if (b >= 0x0A) && (b <= 0x0F) {
		return byte(b - 0x0A + 'A')
	}

	return '0'
}

func char2bin(b byte) byte {

	if (b >= '0') && (b <= '9') {
		return byte(b - '0')
	} else if (b >= 'A') && (b <= 'F') {
		return byte(b - 'A' + 0x0A)
	}

	return 0xFF
}

func HI(b byte) (byte) {
	return (b >> 8) & 0xFF
}

func LOW(b byte) (byte) {
	return b & 0x0F
}

func LRC(data []byte) byte {

	var ucLRC uint8 = 0

	var b byte
	for _, b = range data {
		ucLRC += b
	}

	return uint8(0xFF - ucLRC)
}
