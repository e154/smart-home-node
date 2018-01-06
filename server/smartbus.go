package server

import (
	"time"
	"fmt"
	"github.com/pkg/errors"
	. "github.com/e154/smart-home-node/cache"
)

const (
	ADDRESS uint8 = 0
)

var cache *Cache

type Request struct {
	Line     string        `json: "line"`
	Device   string        `json: "device"`
	Baud     int           `json: "baud"`
	StopBits int           `json: "stop_bits"`
	Sleep    int64         `json: "sleep"`
	Timeout  time.Duration `json: "timeout"`
	Command  []byte        `json: "command"`
	Result   bool          `json: "result"`
}

type Result struct {
	Command   []byte `json: "command"`
	Device    string `json: "device"`
	Result    string `json: "result"`
	Error     string `json: "error"`
	ErrorCode string `json: "error_code"`
}

type Smartbus struct {

}

func (m *Smartbus) Send(request *Request, result *Result) error {

	//fmt.Println("---", request.Command)

	if len(request.Command) == 0 {
		return errors.New("bad command")
	}

	addr := int(request.Command[ADDRESS])
	cacheKey := cache.GetKey(fmt.Sprintf("%d_dev", addr))

	if request.Device == "" {
		if cache.IsExist(cacheKey) {
			request.Device = cache.Get(cacheKey).(string)
		}
	}

	done := make(chan []byte)
	defer func() {
		close(done)
	}()

	var iter int
	Server.ClientQueue(addr,request.Device, func(thread *Thread) {
		if iter > 0 {
			return
		}
		iter++
		//fmt.Println("callback", request.Command)
		err := thread.Send(request, result)
		if err == nil {
			thread.Type = THREAD_SMARTBUS
			cache.Put("node", cacheKey, thread.Dev)
		}

		done <- request.Command
	})

	<- done

	return nil
}

func init() {
	cache = &Cache{
		Cachetime: 3600,
		Name:      "node",
	}
}
