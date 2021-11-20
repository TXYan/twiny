package spider

import (
	"fmt"
	"time"
	"sync"
	"twiny/core"
	"twiny/parser"
	"twiny/scheduler"
	"twiny/downloader"
	"twiny/middleware"
)


type Spider struct {
	// 爬虫基本信息
	name	string
	startUrls	[]*core.Request

	// 下载线程数
	dlThreadNum	int 
	// 解析线程数
	parseThreadNum		int 
	// 默认下载时间间隔
	defaultDlDuration	time.Duration
	// 解析通道
	parseChan			chan *core.Response

	// 等待
	wg					sync.WaitGroup

	// 解析器
	iparser		parser.IParser

	// 定制化调度和下载器
	scheduler	scheduler.IScheduler
	downloader	downloader.IDownloader

	// 中间件
	dlMiddlewares	[]middleware.IDlMiddleware
}

/**
 * 爬虫名称
 */
func (self *Spider) Name() string {
	return self.name
}

/**
 * 爬虫爬取
 */
func (self *Spider) Crawl() {

	// 检查爬虫设置的最小要求 name, startUrls, iparser
	if self.name == "" || len(self.startUrls) == 0 || self.iparser == nil {
		panic("最少需要设置name, startUrls, iparser")
	}

	// 检查scheduler 和 downloader
	if self.defaultDlDuration <= 0 {
		self.defaultDlDuration = 2 * time.Second
	}
	if self.scheduler == nil {
		self.scheduler = scheduler.NewSimpleScheduler(128, self.defaultDlDuration)
	}
	if self.downloader == nil {
		self.downloader = downloader.NewSimpleDownloader(nil)
	}

	// 处理起始urls
	for _, req := range self.startUrls {
		self.scheduler.AddRequest(req)
	}

	// 默认2个线程进行下载
	if self.dlThreadNum == 0 {
		self.dlThreadNum = 2
	}
	for i := 0; i < self.dlThreadNum; i++ {
		self.wg.Add(1)
		go self.handleFireRequest()
	}

	// 默认2个线程进行解析
	if self.parseThreadNum == 0 {
		self.parseThreadNum = 2
	}
	for i := 0; i < self.parseThreadNum; i++ {
		self.wg.Add(1)
		go self.handleResponse()
	}

	defer self.Close()

	self.wg.Wait()
	fmt.Println("spider crawl done")
}

/**
 * 关闭相关通道
 */
func (self *Spider) Close() {
	self.scheduler.Close()
	close(self.parseChan)
}

/**
 * 设置下载线程数
 */
func (self *Spider) DlThreadNum(num int) *Spider {
	self.dlThreadNum = num
	return self
}

/**
 * 设置解析线程数
 */
func (self *Spider) ParseThreadNum(num int) *Spider {
	self.parseThreadNum = num
	return self
}

func (self *Spider) DefaultDlDuration(duration time.Duration) *Spider {
	self.defaultDlDuration = duration
	return self
}

/**
 * 设置调度器 不设置 默认使用SimpleScheduler
 */
func (self *Spider) Scheduler(scheduler scheduler.IScheduler) *Spider {
	self.scheduler = scheduler
	return self
}

/**
 * 设置下载器 不设置 默认使用SimpleDownloader
 */
func (self *Spider) Downloader(downloader downloader.IDownloader) *Spider {
	self.downloader = downloader
	return self
}

/**
 * 设置解析器
 */
func (self *Spider) Parser(parser parser.IParser) *Spider {
	self.iparser = parser
	return self
}

/**
 * 函数方式设置解析器
 */
func (self *Spider) ParserFunc(parseFunc func(resp *core.Response) []*core.Request) *Spider {
	pr := parser.NewParser(parseFunc)
	return self.Parser(pr)
}

/**
 * 增加下载器中间件 按添加顺序执行
 */
func (self *Spider) AddDlMiddleware(dlMiddleware middleware.IDlMiddleware) *Spider {
	self.dlMiddlewares = append(self.dlMiddlewares, dlMiddleware)
	return self
}

/**
 * 函数方式增加下载器中间件 按添加顺序执行
 */
func (self *Spider) AddDlMiddlewareFunc(reqFunc func(req *core.Request) interface{}, respFunc func(resp *core.Response) interface{}) *Spider {
	dlm := middleware.NewDlMiddleware(reqFunc, respFunc)
	return self.AddDlMiddleware(dlm)
}

/**
 * 处理调度器发出的request请求
 */
func (self *Spider) handleFireRequest() {
	defer self.wg.Done()
	
	for req := range self.scheduler.FireRequest() {
		// TODO 中间件合理化
		dlmIdx, newReq, resp := self.dlMiddlwareReqeust(req)
		if newReq != nil {
			// 重新入队列
		 	self.scheduler.AddRequest(newReq)
		 	continue
		}

		// 如果已经有结果则不调用downloader
		if resp == nil {
			resp = self.downloader.Download(req)
		}

		newReq, resp = self.dlMiddlwareResponse(resp, dlmIdx)
		if newReq != nil {
			// 重新入队列
		 	self.scheduler.AddRequest(newReq)
		 	continue
		}

		if resp != nil {
			self.parseChan <- resp
		}
	}
}

/**
 * 处理产生的response 进行解析
 */
func (self *Spider) handleResponse() {
	defer self.wg.Done()

	for resp := range self.parseChan {
		nextReqs := self.iparser.Parse(resp)
		if len(nextReqs) > 0 {
			for _, nextReq := range nextReqs {
				self.scheduler.AddRequest(nextReq)
			}
		}
	}
}

/**
 * 下载器中间件按顺序处理request
 */
func (self *Spider) dlMiddlwareReqeust(req *core.Request) (int, *core.Request, *core.Response) {
	// middlewares
	if len(self.dlMiddlewares) <= 0 {
		return -1, nil, nil
	}


	for idx, middleware := range self.dlMiddlewares {
		result := middleware.ProcessRequest(req)
		/**
		 * 返回 nil 继续下一个中间件
		 * 返回 *core.Request 终止中间件调用 且将返回的request放入到队列中
		 * 返回 *core.Response 终止继续后续的中间件调用 已经激活的按倒序调用
		 */
		 if result != nil {
		 	if newReq, ok := result.(*core.Request); ok {
		 		return idx, newReq, nil
		 	} else if newResp, ok := result.(*core.Response); ok {
		 		return idx, nil, newResp
		 	}
		 }
	}

	return len(self.dlMiddlewares) - 1, nil, nil
}

/**
 * 下载器中间件按执行过ProcessRequest的中间件 倒序执行 ProcessResponse
 */
func (self *Spider) dlMiddlwareResponse(resp *core.Response, maxIdx int) (*core.Request, *core.Response) {

	// middlewares
	if len(self.dlMiddlewares) <= 0 {
		return nil, resp
	}

	/**
	 * 返回 nil 继续按参数resp调用下一个中间件
	 * 返回 *core.Request 终止中间件调用 且将返回的request放入到队列中
	 * 返回 *core.Response 按照返回的response继续调用下一个中间件
	 */
	for idx := maxIdx; idx >= 0; idx-- {
		result := self.dlMiddlewares[idx].ProcessResponse(resp)
		if result != nil {
			if newResp, ok := result.(*core.Response); ok {
				resp = newResp
			} else if newReq, ok := result.(*core.Request); ok {
				return newReq, nil
			}
		}
	}

	return nil, resp
}

/**
 * 创建爬虫实例
 */
func NewSpider(name string, startUrls []*core.Request) *Spider {
	return &Spider {
		name 		: name,
		startUrls	: startUrls,
		parseChan	: make(chan *core.Response),
		dlMiddlewares : make([]middleware.IDlMiddleware, 0),
	}
}