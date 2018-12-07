package command

import (
	"github.com/op/go-logging"
	"encoding/json"
	"github.com/e154/smart-home-node/models/devices"
	"github.com/e154/smart-home-node/common"
)

var (
	log = logging.MustGetLogger("command")
)

type Command struct {
	respFunc       func(data []byte)
	name           string
	args           []string
	requestMessage *common.MessageRequest
}

func NewCommand(respFunc func(data []byte), requestMessage *common.MessageRequest) (command *Command) {

	request := &devices.DevCommandRequest{}
	json.Unmarshal(requestMessage.Command, request)
	command = &Command{
		respFunc:       respFunc,
		name:           request.Name,
		args:           request.Args,
		requestMessage: requestMessage,
	}
	return
}

func (c *Command) Exec() *Command {
	res := ExecuteSync(c.name, c.args...)
	c.response(res)
	return c
}

func (c Command) response(r *Response) {

	respData := &devices.DevCommandResponse{
		Result: r.Out,
		BaseResponse: devices.BaseResponse{
			Error: r.Err,
		},
	}

	data, err := json.Marshal(respData)
	if err != nil {
		log.Error(err.Error())
	}

	response := &common.MessageResponse{
		DeviceId:   c.requestMessage.DeviceId,
		DeviceType: c.requestMessage.DeviceType,
		Properties: c.requestMessage.Properties,
		Response:   data,
		Status:     "success",
	}

	if r.Err != "" || err != nil {
		response.Status = "error"
	}

	responseData, _ := json.Marshal(response)
	c.respFunc(responseData)
}
