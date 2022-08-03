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

package common

import "github.com/e154/smart-home-node/system/serial"

type StatusType string

const (
	Enabled  = StatusType("enabled")
	Disabled = StatusType("disabled")
)

type DeviceType string

const (
	DevTypeModBusRtu = DeviceType("modbus_rtu")
	DevTypeModBusTcp = DeviceType("modbus_tcp")
	DevTypeDefault   = DeviceType("default")
	DevTypeCommand   = DeviceType("command")
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
	Send(entityId string, data interface{})
	EntityId() string
}

type Thread interface {
	//DEPRECATED
	SetParams(baud, timeout, stopBits int) (err error)
	//DEPRECATED
	Open() (err error)
	//DEPRECATED
	Close()
	Restart()
	GetSerial() *serial.Serial
	Device() string
	SetErr()
	SetCon(conn interface{})
	GetCon() interface{}
}
