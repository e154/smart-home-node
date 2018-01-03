package server

import (
	. "github.com/e154/smart-home-node/cache"
	"github.com/e154/smart-home-node/serial"
	. "github.com/e154/smart-home-node/settings"
	"github.com/e154/smart-home-node/serial/smartbus"
	"fmt"
	"errors"
	"encoding/hex"
	"sync"
	"time"
)

const (
	ADDRESS uint8 = 0
)

var cache *Cache

type Request struct {
	Line		string		`json: "line"`
	Device		string		`json: "device"`
	Baud		int		`json: "baud"`
	StopBits	int		`json: "stop_bits"`
	Sleep		int64		`json: "sleep"`
	Timeout		time.Duration	`json: "timeout"`
	Command		[]byte		`json: "command"`
	Result		bool		`json: "result"`
}

type Result struct {
	Command   []byte		`json: "command"`
	Device    string		`json: "device"`
	Result    string		`json: "result"`
	Error     string		`json: "error"`
	ErrorCode string		`json: "error_code"`
}

type Smartbus struct {
	mu	sync.Mutex
}

func (m *Smartbus) Send(request *Request, result *Result) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(request.Command) == 0 {
		return errors.New("command == []")
	}

	conn := &serial.Serial{
		Dev: "",
		Baud: AppConfig.Baud,
		ReadTimeout: AppConfig.Timeout,
		StopBits: AppConfig.StopBits,
	}

	if request.Device != "" {
		conn.Dev = request.Device
	}

	if request.Baud != 0 {
		conn.Baud = request.Baud
	}

	if request.Timeout != 0 {
		conn.ReadTimeout = request.Timeout
	}

	var err error

	if conn.Dev == "" {


		cache_key := cache.GetKey(fmt.Sprintf("%d_dev", request.Command[ADDRESS]))

		//log.Println("send", request.Command)
		//for i := 0; i<5; i++ {

			cache_exist := cache.IsExist(cache_key)
			if cache_exist {
				conn.Dev = cache.Get(cache_key).(string)
				result.Result, err, result.ErrorCode = m.exec(conn, request)
				if err == nil {
					result.Device = conn.Dev
					return nil
				}
			} else {

				//
				devices := serial.DeviceList()
				for _, device := range devices {
					conn.Dev = device
					result.Result, err, result.ErrorCode = m.exec(conn, request)
					if err == nil {
						result.Device = device
						return nil
					}
				}
			}

		//}
	} else {
		for i := 0; i<5; i++ {
			result.Result, err, result.ErrorCode = m.exec(conn, request)
			if err == nil {
				result.Device = conn.Dev
				return nil
			}
		}
	}

	if err != nil {
		result.Error = err.Error()
	}

	return nil
}

func (m *Smartbus) exec(conn *serial.Serial, request *Request) (result string, err error, errcode string) {

	// get cache
	cache_key := cache.GetKey(fmt.Sprintf("%d_dev", request.Command[ADDRESS]))

	if _, err = conn.Open(); err != nil {
		//cache.Delete(cache_key)
		errcode = "SERIAL_PORT_ERROR"
		//log.Printf("error: %s - %s\r\n",conn.Dev, err.Error())
		return
	}
	defer conn.Close()

	modbus := &smartbus.Smartbus{Serial: conn}
	var b []byte
	if b, err = modbus.Send(request.Command); err != nil {
		//cache.Delete(cache_key)
		errcode = "MODBUS_LINE_ERROR"
		//log.Printf("error: %s - %s\r\n",conn.Dev, err.Error())
		return
	}

	result = hex.EncodeToString(b)
	cache.Put("node", cache_key, conn.Dev)

	// bug in the devices need timeout, need fix!!!
	if request.Sleep != 0 {
		time.Sleep(time.Millisecond * time.Duration(request.Sleep))
	}

	return
}

func init() {
	cache = &Cache{
		Cachetime: 3600,
		Name: "node",
	}
}