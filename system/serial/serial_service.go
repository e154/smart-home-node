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
	"io/ioutil"
	"strings"
	"time"

	"github.com/e154/smart-home-node/system/config"
)

type SerialService struct {
	deviceList []string
	serialList []*Serial
	cfg        *config.AppConfig
}

func NewSerialService(cfg *config.AppConfig) *SerialService {
	service := &SerialService{
		deviceList: make([]string, 0),
		serialList: make([]*Serial, 0),
		cfg:        cfg,
	}
	go service.run()
	return service
}

func (s *SerialService) run() {
	for {
		time.Sleep(1 * time.Second)
		s.deviceList = s.DeviceList()
	}
}

func (s *SerialService) DeviceList() []string {

	devices := make([]string, 0)
	contents, _ := ioutil.ReadDir("/dev")

	var found bool
	for _, f := range contents {
		found = false
		for _, serial := range s.cfg.Serial {
			if strings.Contains(f.Name(), serial) {
				if !found {
					devices = append(devices, "/dev/"+f.Name())
					found = true
				}
			}
		}
	}

	return devices
}

func (s *SerialService) SerialList(baud int, readTimeout time.Duration, stopBits int) (serialList []*Serial) {

	serialList = make([]*Serial, 0)

	devList := s.DeviceList()
	if len(devList) == 0 {
		return
	}

	for _, dev := range devList {
		serialPort := &Serial{
			Dev:         dev,
			Baud:        baud,
			ReadTimeout: readTimeout,
			StopBits:    stopBits,
		}

		serialList = append(serialList, serialPort)
	}

	return
}
