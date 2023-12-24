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

package client

import (
	"sync"
	"time"

	"github.com/paulbellamy/ratecounter"
)

type Stat struct {
	sync.Mutex
	min        int64
	max        int64
	latency    int64
	rpsCounter *ratecounter.RateCounter
	avgRequest *ratecounter.AvgRateCounter
	startedAt  time.Time
}

func NewStat() Stat {
	return Stat{
		rpsCounter: ratecounter.NewRateCounter(1 * time.Second),
		avgRequest: ratecounter.NewAvgRateCounter(60 * time.Second),
		startedAt:  time.Now(),
	}
}

func (c *Stat) rpsCounterIncr() {
	c.Lock()
	c.rpsCounter.Incr(1)
	c.Unlock()
}

func (c *Stat) avgStart() time.Time {
	return time.Now()
}

func (c *Stat) avgEnd(startTime time.Time) {
	c.latency = time.Since(startTime).Nanoseconds()

	c.Lock()
	c.avgRequest.Incr(c.latency)

	if c.min == 0 {
		c.min = c.latency
	}

	if c.max == 0 {
		c.max = c.latency
	}

	switch {
	case c.latency < c.min:
		c.min = c.latency
	case c.latency > c.max:
		c.max = c.latency
	}
	c.Unlock()
}

func (c *Stat) GetStat() StateSnapshot {
	c.Lock()
	defer c.Unlock()

	return StateSnapshot{
		Min:       c.min,
		Max:       c.max,
		Latency:   c.latency,
		Rps:       c.rpsCounter.Rate(),
		StartedAt: c.startedAt,
	}
}
