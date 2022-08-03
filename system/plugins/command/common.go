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
	"bytes"
	"os/exec"
	"strings"
)

// IC.Execute "sh", "-c", "echo stdout; echo 1>&2 stderr"
func ExecuteSync(name string, arg ...string) (r *Response) {

	r = &Response{}

	//log.Infof("Execute [SYNC] command: %s %s", name, strings.Trim(fmt.Sprint(arg), "[]"))

	// https://golang.org/pkg/os/exec/#example_Cmd_Start
	cmd := exec.Command(name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		r.Err = err.Error()
		return
	}

	if err := cmd.Wait(); err != nil {
		r.Err = err.Error()
		return
	}

	r.Out = strings.TrimSuffix(stdout.String(), "\n")
	r.Err = strings.TrimSuffix(stderr.String(), "\n")

	return
}
