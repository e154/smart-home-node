package server

import (
	"net"
	"fmt"
	"log"
	"net/rpc"
)

func Start(addr string, port int) {
	server := &Server{
		Addr: addr,
		Port: port,
		clients: make(map[*Client]bool),
	}

	server.Start()
}

type Server struct {
	Addr	string
	Port 	int
	listener net.Listener
	clients map[*Client]bool
}

func (s *Server) Start() (err error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		return
	}

	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}

	rpc.Register(&Modbus{})
	rpc.Register(&Node{})

	log.Printf("Start server on %s:%d\r\n", s.Addr, s.Port)

	go func() {
		for {
			conn, _ := s.listener.Accept()
			go s.AddClient(conn)
		}
	}()

	return
}

func (s *Server) AddClient(conn net.Conn) {

	addr := conn.RemoteAddr().String()
	log.Printf("New client %s\r\n", addr)

	client := &Client{}
	s.clients[client] = true
	defer func() {
		conn.Close()
		delete(s.clients, client)
		log.Printf("Connection %s closed", addr)
	}()

	client.listener(conn)
}
