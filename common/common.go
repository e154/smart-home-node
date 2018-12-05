package common

import (
	"encoding/json"
	"time"
)

type MessageRequest struct {
	DeviceId   int64           `json:"device_id"`
	DeviceType DeviceType      `json:"device_type"`
	Properties json.RawMessage `json:"properties" valid:"Required"`
	Command    json.RawMessage `json:"command"`
}

type MessageResponse struct {
	DeviceId   int64           `json:"device_id"`
	DeviceType DeviceType      `json:"device_type"`
	Properties json.RawMessage `json:"properties"`
	Response   json.RawMessage `json:"response"`
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
