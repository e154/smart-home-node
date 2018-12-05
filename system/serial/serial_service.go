package serial

import (
	"time"
	"github.com/e154/smart-home-node/system/config"
	"io/ioutil"
	"strings"
)

type SerialService struct {
	deviceList []string
	serialList []*Serial
	cfg        *config.AppConfig
}

func NewSerialService(cfg *config.AppConfig) *SerialService {
	service := &SerialService{
		deviceList: make([]string, 0),
		serialList: make([]*Serial, 0),
		cfg:        cfg,
	}
	go service.run()
	return service
}

func (s *SerialService) run() {
	for ; ; {
		time.Sleep(1 * time.Second)
		s.deviceList = s.DeviceList()
	}
}

func (s *SerialService) DeviceList() []string {

	devices := make([]string, 0)
	contents, _ := ioutil.ReadDir("/dev")

	var found bool
	for _, f := range contents {
		found = false
		for _, serial := range s.cfg.Serial {
			if strings.Contains(f.Name(), serial) {
				if !found {
					devices = append(devices, "/dev/"+f.Name())
					found = true
				}
			}
		}
	}

	return devices
}

func (s *SerialService) SerialList(baud int, readTimeout time.Duration, stopBits int) (serialList []*Serial) {

	serialList = make([]*Serial, 0)

	devList := s.DeviceList()
	if len(devList) == 0 {
		return
	}

	for _, dev := range devList {
		serialPort := &Serial{
			Dev:         dev,
			Baud:        baud,
			ReadTimeout: readTimeout,
			StopBits:    stopBits,
		}

		serialList = append(serialList, serialPort)
	}

	return
}
