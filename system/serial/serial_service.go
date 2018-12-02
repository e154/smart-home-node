package serial

import "time"

type SerialService struct {
	deviceList []string
	serialList []*Serial
}

func NewSerialService() *SerialService {
	service := &SerialService{
		deviceList: make([]string, 0),
		serialList: make([]*Serial, 0),
	}
	go service.run()
	return service
}

func (s *SerialService) run() {
	for ;; {
		time.Sleep(1 * time.Second)
		s.deviceList = DeviceList()
	}
}