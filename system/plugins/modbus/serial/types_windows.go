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

//go:build ignore
// +build ignore

package serial

// #include <windows.h>
import "C"

const (
	c_MAXDWORD    = C.MAXDWORD
	c_ONESTOPBIT  = C.ONESTOPBIT
	c_TWOSTOPBITS = C.TWOSTOPBITS
	c_EVENPARITY  = C.EVENPARITY
	c_ODDPARITY   = C.ODDPARITY
	c_NOPARITY    = C.NOPARITY
)

type c_COMMTIMEOUTS C.COMMTIMEOUTS

type c_DCB C.DCB

func toDWORD(val int) C.DWORD {
	return C.DWORD(val)
}

func toBYTE(val int) C.BYTE {
	return C.BYTE(val)
}
