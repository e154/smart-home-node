package server

import (
	"sync"
	. "github.com/e154/smart-home-node/serial"
	. "github.com/e154/smart-home-node/settings"
	"github.com/e154/smart-home-node/serial/smartbus"
	"encoding/hex"
	"errors"
	"time"
	"fmt"
)

type ThreadState string

const (
	THREAD_OK            = ThreadState("ok")
	THREAD_BUSY          = ThreadState("busy")
	THREAD_DEV_NOT_FOUND = ThreadState("port by device not found")
	THREAD_NOT_FOUND     = ThreadState("port not found")
	THREAD_ALL_BUSY      = ThreadState("all ports busy")
)

type ThreadType string

const (
	THREAD_MODBUS   = ThreadType("modbus")
	THREAD_SMARTBUS = ThreadType("smartbus")
)

type Threads map[string]*Thread

func NewThread(dev string) (thread *Thread) {

	serial := &Serial{
		Dev:         dev,
		Baud:        AppConfig.Baud,
		ReadTimeout: AppConfig.Timeout,
		StopBits:    AppConfig.StopBits,
	}

	thread = &Thread{
		Dev:    dev,
		Serial: serial,
		Busy:   false,
	}

	return
}

type Thread struct {
	sync.Mutex
	Busy   bool
	Dev    string
	Client *Client
	Serial *Serial
	Type   ThreadType
}

func (t *Thread) Send(request *Request, result *Result) (err error) {

	if t.Busy {
		err = errors.New("device is busy")
		return
	}

	t.Lock()
	t.Busy = true
	defer func() {
		t.Unlock()
		t.Busy = false
		go Server.ThreadReady(t)
	}()

	//fmt.Println("send ->", request.Command)

	if request.Device != "" {
		t.Serial.Dev = request.Device
	}
	if request.Baud != 0 {
		t.Serial.Baud = request.Baud
	}
	if request.Timeout != 0 {
		t.Serial.ReadTimeout = request.Timeout
	}

	result.Result, err, result.ErrorCode = t.exec(request)
	if err == nil {
		return nil
	}

	if err != nil {
		result.Error = err.Error()
	}

	return
}

func (t *Thread) exec(request *Request) (result string, err error, errcode string) {

	if _, err = t.Serial.Open(); err != nil {
		//cache.Delete(cache_key)
		errcode = "SERIAL_PORT_ERROR"
		//log.Printf("error: %s - %s\r\n",conn.Dev, err.Error())
		return
	}
	defer t.Serial.Close()

	modbus := &smartbus.Smartbus{Serial: t.Serial}
	var b []byte
	if b, err = modbus.Send(request.Command); err != nil {
		//cache.Delete(cache_key)
		errcode = "MODBUS_LINE_ERROR"
		//log.Printf("error: %s - %s\r\n",conn.Dev, err.Error())
		return
	}
	result = hex.EncodeToString(b)

	// bug in the devices need timeout, need fix!!!
	if request.Sleep != 0 {
		time.Sleep(time.Millisecond * time.Duration(request.Sleep))
	}

	return
}

func (t *Thread) Remove() {
	fmt.Println("Remove device", t.Dev)
}
