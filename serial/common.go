package serial

import (
	"io/ioutil"
	"strings"
	"time"
)

const (
	BAUD = 19200
	READ_TIMEOUT = time.Millisecond * 200
	STOP_BITS = 2
)

func DeviceList() []string {

	devices := []string{}
	contents, _ := ioutil.ReadDir("/dev")

	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "cu.SLAB_USB") ||
			strings.Contains(f.Name(), "ttyS") ||
			strings.Contains(f.Name(), "ttyUSB") {
			devices = append(devices, "/dev/" + f.Name())
		}
	}

	return devices
}

func SerialList() (serial_list []*Serial) {

	dev_list := DeviceList()
	if len(dev_list) == 0 {
		return
	}

	for _, dev := range dev_list {
		serial_port := &Serial{
			Dev: dev,
			Baud: BAUD,
			ReadTimeout: READ_TIMEOUT,
			StopBits: STOP_BITS,
		}

		serial_list = append(serial_list, serial_port)
	}

	return
}