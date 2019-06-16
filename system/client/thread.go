package client

import (
	"sync"
	"time"
	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/system/serial"
)

type Thread struct {
	sync.Mutex
	Busy      bool
	Dev       string
	serialDev *serial.Serial
	baud      int
	timeout   int64
	stopBits  int
	errors    int
	conn      interface{}
	Active    bool
	blockList []int64
}

type Threads map[string]*Thread

func NewThread(dev string) (thread *Thread) {

	thread = &Thread{
		Dev:       dev,
		Busy:      false,
		Active:    true,
		blockList: make([]int64, 0),
	}

	return
}

func (t *Thread) Exec(cb common.ThreadCaller) (resp *common.MessageResponse, err error) {

	t.Lock()
	t.Busy = true
	defer func() {
		t.Unlock()
		t.Busy = false
	}()

	resp, err = cb.Exec(t)

	return
}

func (t *Thread) SetParams(baud, timeout, stopBits int) (err error) {

	var restart bool
	if t.baud != baud {
		t.baud = baud
		restart = true
	}
	if t.timeout != int64(timeout) {
		t.timeout = int64(timeout)
		restart = true
	}
	if stopBits != stopBits {
		stopBits = stopBits
		restart = true
	}
	if restart {
		t.Restart()
	}
	return
}

func (t *Thread) Open() (err error) {

	log.Warningf("open device %s", t.Dev)

	t.serialDev = &serial.Serial{
		Dev:         t.Dev,
		Baud:        t.baud,
		ReadTimeout: time.Duration(t.timeout),
		StopBits:    t.stopBits,
	}

	if _, err = t.serialDev.Open(); err != nil {
		log.Warningf("%s - %s\r\n", t.Dev, err.Error())
		return
	}

	t.errors = 0

	return
}

func (t *Thread) Close() {
	if t.serialDev != nil {
		log.Warningf("close device %s", t.Dev)
		t.serialDev.Close()
	}
	return
}

//DEPRECATED
func (t *Thread) Restart() {
	t.Close()
	t.Open()
}

func (t *Thread) GetSerial() *serial.Serial {
	return t.serialDev
}

func (t *Thread) Device() string {
	return t.Dev
}

func (t *Thread) SetErr() {
	t.errors++
	if t.errors > 30 {
		t.Restart()
	}
}

func (t *Thread) SetCon(conn interface{}) {
	t.conn = conn
}

func (t *Thread) GetCon() interface{} {
	return t.conn
}

func (t *Thread) Disable() {
	t.Active = false
	t.errors = 0
	t.conn = nil
	t.blockList = make([]int64, 0)
}

func (t *Thread) Enable() {
	t.Active = true
}