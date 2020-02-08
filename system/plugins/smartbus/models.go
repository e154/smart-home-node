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

package smartbus

import "time"

type Request struct {
	Line     string        `json:"line"`
	Device   string        `json:"device"`
	Baud     int           `json:"baud"`
	StopBits int           `json:"stop_bits"`
	Sleep    int64         `json:"sleep"`
	Timeout  time.Duration `json:"timeout"`
	Command  []byte        `json:"command"`
	Result   bool          `json:"result"`
}

type Result struct {
	Command   []byte `json:"command"`
	Device    string `json:"device"`
	Result    string `json:"result"`
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}
