package server

import (
	"net"
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"time"
	"github.com/e154/smart-home-node/serial"
)

const (
	TIME_TICK = 30 * time.Second
)

var (
	Server *server
)

type Queue struct {
	Dev		string
	Cb		func(thread *Thread)
}

func Start(addr string, port int) {
	Server = &server{
		Addr:       addr,
		Port:       port,
		lastUpdate: time.Now(),
		clients:    make(map[*Client]bool),
		pool:       make(map[string]*Thread),
	}

	Server.UpdateThreads()
	Server.run()
}

type server struct {
	Addr       string
	Port       int
	listener   net.Listener
	clients    map[*Client]bool
	sync.RWMutex
	lastUpdate time.Time
	deviceList []string
	pool       Threads
	queue	   sync.Map
}

func (s *server) run() (err error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		return
	}

	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}

	rpc.Register(&Smartbus{})
	rpc.Register(&Node{})

	log.Printf("start server on %s:%d\r\n", s.Addr, s.Port)

	go func() {
		for {
			conn, _ := s.listener.Accept()
			go s.AddClient(conn)
		}
	}()

	go func() {
		ticker := time.Tick(TIME_TICK)
		for {
			select {
			case <-ticker:
				s.UpdateThreads()
				s.UpdateCLientQueue()
			}
		}
	}()

	return
}

func (s *server) AddClient(conn net.Conn) {

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

// обновление пула активных устройств
//
func (s *server) UpdateThreads() {

	s.Lock()
	defer s.Unlock()

	s.deviceList = serial.DeviceList()

	//remove threads
	for pKey, thread := range s.pool {
		var exist bool
		for _, dev := range s.deviceList {
			if dev == pKey {
				exist = true
			}
		}
		if exist {
			continue
		}
		thread.Remove()
		fmt.Println("Remove device from pool:", pKey)
		delete(s.pool, pKey)
	}

	for _, dev := range s.deviceList {
		if _, ok := s.pool[dev]; ok {
			continue
		}

		fmt.Println("Add device to pool:", dev)
		s.pool[dev] = NewThread(dev)
	}
}

func (s *server) GetThread(dev string) (*Thread, ThreadState) {

	s.Lock()
	defer s.Unlock()

	//если нет портов
	if len(s.pool) == 0 {
		return nil, THREAD_DEV_NOT_FOUND
	}

	//выбор порта по названию
	if dev != "" {
		if thread, ok := s.pool[dev]; ok {
			if thread.Busy {
				return nil, THREAD_BUSY
			}

			return thread, THREAD_OK
		}
		//по названию не найдено
		return nil, THREAD_DEV_NOT_FOUND
	}

	//выбор свободного порта
	for _, thread := range s.pool {
		if thread.Busy {
			continue
		}

		return thread, THREAD_OK
	}

	//если все порты заняты
	if len(s.pool) > 0 {
		return nil, THREAD_ALL_BUSY
	}

	return nil, THREAD_NOT_FOUND
}

func (s *server) ClientQueue(addr int, dev string,f func(thread *Thread)) {

	s.queue.Store(addr, &Queue{dev, f})

	s.UpdateCLientQueue()
}

func (s *server) UpdateCLientQueue() {

	s.queue.Range(func(key, value interface{}) bool {

		cli := value.(*Queue)

		thread, status := s.GetThread(cli.Dev)
		if status == THREAD_OK {
			go cli.Cb(thread)
			s.queue.Delete(key)
		}

		return true
	})
}

func (s *server) ThreadReady(thread *Thread) {
	s.UpdateCLientQueue()
}