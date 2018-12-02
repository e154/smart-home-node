package serial

import (
	"io/ioutil"
	"strings"
	"time"
)

func DeviceList() []string {

	devices := make([]string, 0)

	contents, _ := ioutil.ReadDir("/dev")

	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "tty.SLAB_USB") ||
			strings.Contains(f.Name(), "ttyS") ||
			strings.Contains(f.Name(), "ttyUSB") {
			devices = append(devices, "/dev/"+f.Name())
		}
	}

	return devices
}

func SerialList(baud int, readTimeout time.Duration, stopBits int) (serialList []*Serial) {

	serialList = make([]*Serial, 0)

	devList := DeviceList()
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
