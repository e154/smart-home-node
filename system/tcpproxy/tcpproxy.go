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
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"github.com/e154/smart-home-node/system/config"
)

type TcpProxy struct {
	cfg     *config.AppConfig
	quit    chan struct{}
	ln      *net.TCPListener
	mx      sync.Mutex
	clients map[string]*Proxy
}

func NewTcpProxy(lc fx.Lifecycle, cfg *config.AppConfig) *TcpProxy {
	proxy := &TcpProxy{
		cfg:     cfg,
		quit:    make(chan struct{}),
		mx:      sync.Mutex{},
		clients: make(map[string]*Proxy),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return proxy.Start()
		},
		OnStop: func(ctx context.Context) (err error) {
			return proxy.Shutdown()
		},
	})

	return proxy
}

func (p *TcpProxy) Start() error {
	go p.runServer()
	return nil
}

func (p *TcpProxy) Shutdown() error {
	log.Info("Shutdown")

	if p.ln == nil {
		return nil
	}

	if err := p.ln.Close(); err != nil {
		log.Error(err.Error())
	}

	p.quit <- struct{}{}

	count := len(p.clients)
	if count == 0 {
		return nil
	}

	log.Infof("total clients %d", count)

	for _, cli := range p.clients {
		if cli != nil {
			cli.Stop()
		}
	}
	return nil
}

func (p *TcpProxy) runServer() {

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", p.cfg.ProxyPort))
	if err != nil {
		log.Warnf("Failed to resolve local address: %s", err)
		return
	}

	p.ln, err = net.ListenTCP("tcp", laddr)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			return
		}
		log.Warnf("Failed to open local port to listen: %s", err)
		return
	}

	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.cfg.MqttIp, p.cfg.MqttPort))
	if err != nil {
		log.Warnf("Failed to resolve remote address: %s", err)
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
			log.Warnf("Failed to accept connection '%s'", err)
			continue
		}

		go p.addClient(conn, laddr, raddr)
	}
}

func (p *TcpProxy) addClient(conn *net.TCPConn, laddr, raddr *net.TCPAddr) {
	id := uuid.NewString()

	var pr *Proxy
	pr = New(conn, laddr, raddr, id)

	p.clients[id] = pr
	pr.Start()

	if _, ok := p.clients[id]; ok {
		delete(p.clients, id)
	}
}
