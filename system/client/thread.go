package client

import (
	"sync"
	"github.com/e154/smart-home-node/common"
)

type Thread struct {
	sync.Mutex
	Busy bool
	Dev  string
}

type Threads map[string]*Thread

func NewThread(dev string) (thread *Thread) {

	thread = &Thread{
		Dev:  dev,
		Busy: false,
	}

	return
}

func (t *Thread) Send(cb common.ThreadCaller) (resp *common.MessageResponse, err error) {

	t.Lock()
	t.Busy = true
	defer func() {
		t.Unlock()
		t.Busy = false
	}()

	resp, err = cb.Exec(t.Dev)

	return
}
