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

// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs types_windows.go

package serial

const (
	c_MAXDWORD    = 0xffffffff
	c_ONESTOPBIT  = 0x0
	c_TWOSTOPBITS = 0x2
	c_EVENPARITY  = 0x2
	c_ODDPARITY   = 0x1
	c_NOPARITY    = 0x0
)

type c_COMMTIMEOUTS struct {
	ReadIntervalTimeout         uint32
	ReadTotalTimeoutMultiplier  uint32
	ReadTotalTimeoutConstant    uint32
	WriteTotalTimeoutMultiplier uint32
	WriteTotalTimeoutConstant   uint32
}

type c_DCB struct {
	DCBlength  uint32
	BaudRate   uint32
	Pad_cgo_0  [4]byte
	WReserved  uint16
	XonLim     uint16
	XoffLim    uint16
	ByteSize   uint8
	Parity     uint8
	StopBits   uint8
	XonChar    int8
	XoffChar   int8
	ErrorChar  int8
	EofChar    int8
	EvtChar    int8
	WReserved1 uint16
}

func toDWORD(val int) uint32 {
	return uint32(val)
}

func toBYTE(val int) uint8 {
	return uint8(val)
}
