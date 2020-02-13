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

package main

import (
	"fmt"
	"github.com/e154/smart-home-node/system/client"
	"github.com/e154/smart-home-node/system/graceful_service"
	l "github.com/e154/smart-home-node/system/logging"
	"github.com/e154/smart-home-node/system/tcpproxy"
	"github.com/e154/smart-home-node/version"
	"github.com/op/go-logging"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	log = logging.MustGetLogger("main")
)

func main() {

	args := os.Args[1:]
	for _, arg := range args {
		switch arg {
		case "-v", "--version":
			fmt.Printf(version.ShortVersionBanner, version.GetHumanVersion())
			return
		default:
			fmt.Printf(version.VerboseVersionBanner, "v2", os.Args[0])
			return
		}
	}

	start()
}

func start() {

	fmt.Printf(version.ShortVersionBanner, "")

	container := BuildContainer()
	err := container.Invoke(func(
		graceful *graceful_service.GracefulService,
		lx *logrus.Logger,
		client *client.Client,
		server *tcpproxy.TcpProxy) {

		l.Initialize(lx)
		client.Connect()
		go server.Start()

		graceful.Wait()
	})

	if err != nil {
		panic(err.Error())
	}
}
