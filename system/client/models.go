package client

import (
	"encoding/json"
	. "github.com/e154/smart-home-node/common"
	"time"
)

type MessageReq struct {
	DeviceId   int64           `json:"device_id"`
	DeviceType DeviceType      `json:"device_type"`
	Properties json.RawMessage `json:"properties" valid:"Required"`
	Command    []byte          `json:"command"`
}

type MessageResp struct {
	DeviceId   int64           `json:"device_id"`
	DeviceType DeviceType      `json:"device_type"`
	Properties json.RawMessage `json:"properties"`
	Response   []byte          `json:"response"`
	Status     string          `json:"status"`
}

type ClientStatusModel struct {
	Status    ClientStatus `json:"status"`
	Thread    int          `json:"thread"`
	Rps       int64        `json:"rps"`
	Min       int64        `json:"min"`
	Max       int64        `json:"max"`
	StartedAt time.Time    `json:"started_at"`
}
