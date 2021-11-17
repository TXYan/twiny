package scheduler

import (
	"time"
	"twiny/core"
)

type SimpleScheduler struct {
	incomeChan chan *core.Request
	outcomeChan chan *core.Request 
	duration	time.Duration
}

func (self *SimpleScheduler) AddRequest(req *core.Request) {
	self.incomeChan <- req
}

func (self *SimpleScheduler) FireRequest() <- chan *core.Request {
	return self.outcomeChan;
}

func (self *SimpleScheduler) Close() {
	close(self.incomeChan)
	close(self.outcomeChan)
}

func (self *SimpleScheduler) start() {
	ticker := time.NewTicker(self.duration)
	for {
		select {
		case <- ticker.C:
			self.outcomeChan <- <- self.incomeChan
		}
	}
}

func NewSimpleScheduler(queuCap int, duration time.Duration) *SimpleScheduler {
	scheduler := &SimpleScheduler {
		incomeChan : make(chan *core.Request, queuCap),
		outcomeChan : make(chan *core.Request),
		duration    : duration,
	}

	go scheduler.start()

	return scheduler
}