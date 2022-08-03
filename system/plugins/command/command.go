// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package command

import (
	"encoding/json"

	"github.com/e154/smart-home-node/common"
	"github.com/e154/smart-home-node/common/logger"
	"github.com/e154/smart-home-node/models/devices"
)

var (
	log = logger.MustGetLogger("command")
)

type Command struct {
	respFunc       func(entityId string, data []byte)
	name           string
	args           []string
	requestMessage *common.MessageRequest
}

func NewCommand(respFunc func(entityId string, data []byte), requestMessage *common.MessageRequest) (command *Command) {

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
		EntityId:   c.requestMessage.EntityId,
		DeviceType: c.requestMessage.DeviceType,
		Properties: c.requestMessage.Properties,
		Response:   data,
		Status:     "success",
	}

	if r.Err != "" || err != nil {
		response.Status = "error"
	}

	responseData, _ := json.Marshal(response)
	c.respFunc(c.requestMessage.EntityId, responseData)
}
