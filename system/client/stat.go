package client

import (
	"github.com/paulbellamy/ratecounter"
	"time"
)

type Stat struct {
	min        int64
	max        int64
	rpsCounter *ratecounter.RateCounter
	avgRequest *ratecounter.AvgRateCounter
	startTime  time.Time
}

func (c *Stat) rpsCounterIncr() {
	c.rpsCounter.Incr(1)
}

func (c *Stat) avgStart() {
	c.startTime = time.Now()
}

func (c *Stat) avgEnd() {
	total := time.Since(c.startTime).Nanoseconds()
	c.avgRequest.Incr(total)

	switch {
	case total < c.min || c.min == 0:
		c.min = total
	case total > c.max || c.max == 0:
		c.max = total
	}
}