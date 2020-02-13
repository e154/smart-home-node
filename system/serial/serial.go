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

package serial

import (
	"github.com/e154/serial"
	"time"
)

type Serial struct {
	Dev				string
	Baud			int
	ReadTimeout 	time.Duration
	StopBits 		int
	config			*serial.Config
	Port			*serial.Port
}

func (s *Serial) Open() (*Serial, error) {

	s.config = &serial.Config{
		Name: s.Dev,
		Baud: s.Baud,
		StopBits: serial.StopBits(s.StopBits),
		ReadTimeout: s.ReadTimeout,
	}

	var err error
	if s.Port, err = serial.OpenPort(s.config); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Serial) Close() (*Serial, error) {

	if s.Port != nil {
		return s, s.Port.Close()
	}

	return s, nil
}
