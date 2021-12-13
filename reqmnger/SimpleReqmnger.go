package reqmnger

import (
	"twiny/core"
)

type SimpleReqMnger struct {
	reqchan		chan *core.Request
}


func NewSimpleReqMnger(size int) *SimpleReqMnger {
	return &SimpleReqMnger {
		reqchan		:	make(chan *core.Request, size),
	}
}

func (self *SimpleReqMnger) AddReqs(reqs ...*core.Request) {
	for _, req := range reqs {
		self.reqchan <- req
	}
}

func (self *SimpleReqMnger) NextReq() *core.Request {
	return <- self.reqchan
}