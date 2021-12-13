package scheduler

import (
	"time"
	"twiny/core"
	"twiny/reqmnger"
	"twiny/downloader"
	"twiny/config"
)

type Scheduler struct {
	reqmanager 	reqmnger.Ireqmnger
	scheConfig 	config.SchedulerConfig
	downloader	*downloader.Downloader

	fireChan	chan *core.Request
	parseChan	chan *core.Response

}

func NewScheduler(scheConfig config.SchedulerConfig, dler *downloader.Downloader) *Scheduler {
	return &Scheduler {
		scheConfig	: scheConfig,
		reqmanager	: scheConfig.GetReqManager(),
		downloader	: dler,
		fireChan	: make(chan *core.Request, scheConfig.GetDlThreadNum()),
		parseChan	: make(chan *core.Response, scheConfig.GetParseThreadNum()),
	}
}

func (self *Scheduler) Start() {
	reqs := self.scheConfig.GetStartReqs()
	if len(reqs) == 0 {
		panic("start requests is empty")
	}

	if self.reqmanager == nil {
		self.reqmanager = reqmnger.NewSimpleReqMnger(2 * self.scheConfig.GetDlThreadNum())
	}

	for i := 0; i < self.scheConfig.GetParseThreadNum(); i++ {
		go self.parseResponse()
	}

	for i := 0; i < self.scheConfig.GetDlThreadNum(); i++ {
		go self.dlRequest()
	}

	// 将其实url加入管理器
	self.reqmanager.AddReqs(reqs...)
	self.fireRequest()
}

func (self *Scheduler) Close() {
	close(self.fireChan)
	close(self.parseChan)
}

func (self *Scheduler) fireRequest() {
	ticker := time.NewTicker(self.scheConfig.GetDlDuration())
	for {
		<- ticker.C
		req := self.reqmanager.NextReq()
		self.fireChan <- req
	}
}

func (self *Scheduler) dlRequest() {
	for {
		select {
		case req := <- self.fireChan:
			// 中间件处理
			newReq, resp := self.downloader.Download(req)
			if newReq != nil {
				self.reqmanager.AddReqs(newReq)
			}
			if resp != nil {
				self.parseChan <- resp
			}
		}
	}
}

func (self *Scheduler) parseResponse() {
	for {
		select {
		case resp := <- self.parseChan:
			newReqs := self.scheConfig.GetParser().Parse(resp)
			if newReqs != nil && len(newReqs) > 0 {
				self.reqmanager.AddReqs(newReqs...)
			}
		}
	}
}

