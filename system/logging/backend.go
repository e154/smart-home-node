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

package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/op/go-logging"
)

type LogBackend struct {
	L *logrus.Logger
}

func NewLogBackend(logger *logrus.Logger) *LogBackend {
	return &LogBackend{L: logger}
}

func (b *LogBackend) Log(level logging.Level, calldepth int, rec *logging.Record) error {

	s := rec.Formatted(calldepth + 1)
	switch level {
	case logging.CRITICAL:
		b.L.Level = logrus.FatalLevel
		b.L.Fatal(s)
	case logging.ERROR:
		b.L.Level = logrus.ErrorLevel
		b.L.Error(s)
	case logging.WARNING:
		b.L.Level = logrus.WarnLevel
		b.L.Warning(s)
	case logging.INFO, logging.NOTICE:
		b.L.Level = logrus.InfoLevel
		b.L.Info(s)
	case logging.DEBUG:
		b.L.Level = logrus.DebugLevel
		b.L.Debug(s)
	}
	return nil
}
