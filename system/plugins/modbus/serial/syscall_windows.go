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

// +build ignore

package serial

// Windows api calls

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output zsyscall_windows.go $GOFILE

//sys GetCommState(handle syscall.Handle, dcb *c_DCB) (err error)
//sys SetCommState(handle syscall.Handle, dcb *c_DCB) (err error)
//sys GetCommTimeouts(handle syscall.Handle, timeouts *c_COMMTIMEOUTS) (err error)
//sys SetCommTimeouts(handle syscall.Handle, timeouts *c_COMMTIMEOUTS) (err error)
