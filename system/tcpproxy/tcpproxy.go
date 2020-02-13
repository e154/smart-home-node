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

package tcpproxy

import (
	"fmt"
	"github.com/e154/smart-home-node/system/config"
	"github.com/e154/smart-home-node/system/graceful_service"
	"github.com/e154/smart-home-node/system/uuid"
	"net"
	"strings"
	"sync"
	"time"
)

type TcpProxy struct {
	cfg     *config.AppConfig
	quit    chan struct{}
	ln      *net.TCPListener
	mx      sync.Mutex
	clients map[string]*Proxy
}

func NewTcpProxy(cfg *config.AppConfig,
	graceful *graceful_service.GracefulService) *TcpProxy {
	proxy := &TcpProxy{
		cfg:     cfg,
		quit:    make(chan struct{}),
		mx:      sync.Mutex{},
		clients: make(map[string]*Proxy),
	}
	graceful.Subscribe(proxy)
	return proxy
}

func (p *TcpProxy) Start() {
	p.runServer()
}

func (p *TcpProxy) Shutdown() {
	log.Info("Shutdown")

	if p.ln == nil {
		return
	}

	if err := p.ln.Close(); err != nil {
		log.Error(err.Error())
	}

	p.quit <- struct{}{}

	count := len(p.clients)
	if count == 0 {
		return
	}

	log.Infof("total clients %d", count)

	for _, cli := range p.clients {
		if cli != nil {
			cli.Stop()
		}
	}
}

func (p *TcpProxy) runServer() {

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", p.cfg.ProxyPort))
	if err != nil {
		log.Warningf("Failed to resolve local address: %s", err)
		return
	}

	p.ln, err = net.ListenTCP("tcp", laddr)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			return
		}
		log.Warningf("Failed to open local port to listen: %s", err)
		return
	}

	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.cfg.MqttIp, p.cfg.MqttPort))
	if err != nil {
		log.Warningf("Failed to resolve remote address: %s", err)
		return
	}

	log.Infof("Serving server at tcp://[::]:%d", p.cfg.ProxyPort)

	for {
		select {
		case <-p.quit:
			log.Info("Connection closed...")
			return
		default:
			time.Sleep(time.Second)
		}

		conn, err := p.ln.AcceptTCP()
		if err != nil {
			log.Warningf("Failed to accept connection '%s'", err)
			continue
		}

		go p.addClient(conn, laddr, raddr)
	}
}

func (p *TcpProxy) addClient(conn *net.TCPConn, laddr, raddr *net.TCPAddr) {
	id := uuid.NewV4().String()

	var pr *Proxy
	pr = New(conn, laddr, raddr, id)

	p.clients[id] = pr
	pr.Start()

	if _, ok := p.clients[id]; ok {
		delete(p.clients, id)
	}
}
