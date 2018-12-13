package common

import "github.com/e154/smart-home-node/system/serial"

type StatusType string

const (
	Enabled  = StatusType("enabled")
	Disabled = StatusType("disabled")
)

type DeviceType string

const (
	DevTypeSmartBus = DeviceType("smartbus")
	DevTypeModBus   = DeviceType("modbus")
	DevTypeZigbee   = DeviceType("zigbee")
	DevTypeDefault  = DeviceType("default")
	DevTypeCommand  = DeviceType("command")
)

type ThreadState string

const (
	ThreadOk          = ThreadState("ok")
	ThreadBusy        = ThreadState("busy")
	ThreadDevNotFound = ThreadState("port by device not found")
	ThreadNotFound    = ThreadState("port not found")
	ThreadAllBusy     = ThreadState("all ports busy")
)

type ClientStatus string

const (
	StatusEnabled  = ClientStatus("enabled")
	StatusDisabled = ClientStatus("disabled")
	StatusError    = ClientStatus("error")
	StatusBusy     = ClientStatus("busy")
)

type ThreadCaller interface {
	Exec(t Thread) (resp *MessageResponse, err error)
	Send(data interface{})
	DeviceId() int64
}

type Thread interface {
	SetParams(baud, timeout, stopBits int) (err error)
	Open() (err error)
	Close()
	Restart()
	GetSerial() *serial.Serial
	Device() string
	SetErr()
}