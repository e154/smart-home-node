package client

import (
	"encoding/json"
	. "github.com/e154/smart-home-node/common"
)

type Message struct {
	DeviceId   int64           `json:"device_id"`
	DeviceType DeviceType      `json:"device_type"`
	Properties json.RawMessage `json:"properties" valid:"Required"`
}
