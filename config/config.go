package config


import (
	"time"
	"net/http"
	"twiny/core"
	"twiny/parser"
	"twiny/reqmnger"
	"twiny/middleware"
)

type SchedulerConfig struct {
	// 下载线程数
	dlThreadNum		int 
	// 解析线程数
	parseThreadNum	int 
	// 下载间隔
	dlDuration		time.Duration
	// 请求管理器
	reqmanager 		reqmnger.Ireqmnger
	startReqs		[]*core.Request
	iparser			parser.IParser
}


/**
 * 设置下载线程数
 */
func (self *SchedulerConfig) DlThreadNum(num int) {
	self.dlThreadNum = num
}

func (self *SchedulerConfig) GetDlThreadNum() int {
	return self.dlThreadNum
}

/**
 * 设置解析线程数
 */
func (self *SchedulerConfig) ParseThreadNum(num int) {
	self.parseThreadNum = num
}

func (self *SchedulerConfig) GetParseThreadNum() int {
	return self.parseThreadNum
}

func (self *SchedulerConfig) DlDuration(duration time.Duration) {
	self.dlDuration = duration
}

func (self *SchedulerConfig) GetDlDuration() time.Duration {
	return self.dlDuration
}

func (self *SchedulerConfig) ReqManager(reqMnger reqmnger.Ireqmnger) {
	self.reqmanager = reqMnger
}

func (self *SchedulerConfig) GetReqManager() reqmnger.Ireqmnger {
	return self.reqmanager
}

func (self *SchedulerConfig) StartReqs(reqs []*core.Request) *SchedulerConfig {
	self.startReqs = reqs
	return self
}

func (self *SchedulerConfig) GetStartReqs() []*core.Request {
	return self.startReqs
}

func (self *SchedulerConfig) Parser(iparser parser.IParser) {
	self.iparser = iparser
}

func (self *SchedulerConfig) GetParser() parser.IParser {
	return self.iparser
}




type DownloaderConfig struct {
	dlMiddlewares	[]middleware.IDlMiddleware
	dlClient		*http.Client
}

func (self *DownloaderConfig) AddDlMiddleware(dlMiddleware middleware.IDlMiddleware) {
	self.dlMiddlewares = append(self.dlMiddlewares, dlMiddleware)
}

func (self *DownloaderConfig) GetDlMiddlewares() []middleware.IDlMiddleware {
	return self.dlMiddlewares
}


func (self *DownloaderConfig) DlClient(client *http.Client) {
	self.dlClient = client
}

func (self *DownloaderConfig) GetDlClient() *http.Client {
	return self.dlClient
}