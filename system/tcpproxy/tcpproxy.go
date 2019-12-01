package tcpproxy

import (
	"fmt"
	"github.com/e154/smart-home-gate/system/uuid"
	"github.com/e154/smart-home-node/system/config"
	"github.com/e154/smart-home-node/system/graceful_service"
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
	go proxy.runServer()
	graceful.Subscribe(proxy)
	return proxy
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

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", p.cfg.MqttPort))
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

	log.Infof("Serving server at tcp://[::]:%d", p.cfg.MqttPort)

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