package tcpproxy

import (
	"fmt"
	"github.com/e154/smart-home-node/system/config"
	"github.com/e154/smart-home-node/system/graceful_service"
	"net"
	"regexp"
	"strings"
)

var (
	matchid = uint64(0)
	connid  = uint64(0)
)

type TcpProxy struct {
	proxy *Proxy
	cfg   *config.AppConfig
	quit  chan struct{}
}

func NewTcpProxy(cfg *config.AppConfig,
	graceful *graceful_service.GracefulService) *TcpProxy {
	proxy := &TcpProxy{
		cfg:  cfg,
		quit: make(chan struct{}),
	}
	go proxy.runServer()
	graceful.Subscribe(proxy)
	return proxy
}

func (p *TcpProxy) Shutdown() {
	p.quit <- struct{}{}
}

func (p *TcpProxy) runServer() {
	log.Infof("Serving server at tcp://[::]:%d", p.cfg.ProxyPort)

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", p.cfg.ProxyPort))
	if err != nil {
		log.Warningf("Failed to resolve local address: %s", err)
		return
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Warningf("Failed to open local port to listen: %s", err)
		return
	}

	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.cfg.MqttIp, p.cfg.MqttPort))
	if err != nil {
		log.Warningf("Failed to resolve remote address: %s", err)
		return
	}

	matcher := createMatcher("")
	replacer := createReplacer("")

	for {
		select {
		case <-p.quit:
			return
		default:

		}

		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Warningf("Failed to accept connection '%s'", err)
			continue
		}

		var pr *Proxy
		pr = New(conn, laddr, raddr)

		pr.Matcher = matcher
		pr.Replacer = replacer

		go pr.Start()
	}
}

func createMatcher(match string) func([]byte) {
	if match == "" {
		return nil
	}
	re, err := regexp.Compile(match)
	if err != nil {
		log.Warningf("Invalid match regex: %s", err)
		return nil
	}

	log.Infof("Matching %s", re.String())
	return func(input []byte) {
		ms := re.FindAll(input, -1)
		for _, m := range ms {
			matchid++
			log.Infof("Match #%d: %s", matchid, string(m))
		}
	}
}

func createReplacer(replace string) func([]byte) []byte {
	if replace == "" {
		return nil
	}
	//split by / (TODO: allow slash escapes)
	parts := strings.Split(replace, "~")
	if len(parts) != 2 {
		log.Warningf("Invalid replace option")
		return nil
	}

	re, err := regexp.Compile(string(parts[0]))
	if err != nil {
		log.Warningf("Invalid replace regex: %s", err)
		return nil
	}

	repl := []byte(parts[1])

	log.Infof("Replacing %s with %s", re.String(), repl)
	return func(input []byte) []byte {
		return re.ReplaceAll(input, repl)
	}
}
