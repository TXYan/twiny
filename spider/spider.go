package spider

import (
	"fmt"
	"time"
	"runtime"
	"github.com/TXYan/twiny/core"
	"github.com/TXYan/twiny/config"
	"github.com/TXYan/twiny/parser"
	"github.com/TXYan/twiny/scheduler"
	"github.com/TXYan/twiny/middleware"
	"github.com/TXYan/twiny/reqmnger"
	"github.com/TXYan/twiny/downloader"
)

type SchedulerConfig = config.SchedulerConfig
type DownloaderConfig = config.DownloaderConfig

/**
 * URL 管理器
 * 网页下载器
 * 爬虫调度器
 * 
 * 网页解析器
 * 数据处理器
 */

type Spider struct {
	*SchedulerConfig
	*DownloaderConfig
	// 爬虫基本信息
	name			string
	stop			chan bool
}

/**
 * 爬虫爬取
 */
func (self *Spider) Crawl() {

	// 检查爬虫设置的最小要求 name, startUrls, iparser
	if self.name == "" || len(self.GetStartReqs()) == 0 || self.GetParser() == nil {
		panic("最少需要设置name, startUrls, iparser")
	}

	dler := downloader.NewDownloader(*self.DownloaderConfig)
	sch := scheduler.NewScheduler(*self.SchedulerConfig, dler)

	defer sch.Close()
	go sch.Start()

	<- self.stop

	fmt.Println("spider crawl done")
}

func (self *Spider) Close() {
	close(self.stop)
}

/**
 * 爬虫名称
 */
func (self *Spider) Name() string {
	return self.name
}

func (self *Spider) StartReqs(reqs []*core.Request) *Spider {
	self.SchedulerConfig.StartReqs(reqs)
	return self
}

/**
 * 设置下载线程数
 */
func (self *Spider) DlThreadNum(num int) *Spider {
	self.SchedulerConfig.DlThreadNum(num)
	return self
}

/**
 * 设置解析线程数
 */
func (self *Spider) ParseThreadNum(num int) *Spider {
	self.SchedulerConfig.ParseThreadNum(num)
	return self
}

func (self *Spider) DlDuration(duration time.Duration) *Spider {
	self.SchedulerConfig.DlDuration(duration)
	return self
}

/**
 * 设置解析器
 */
func (self *Spider) Parser(parser parser.IParser) *Spider {
	self.SchedulerConfig.Parser(parser)
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
	self.DownloaderConfig.AddDlMiddleware(dlMiddleware)
	return self
}

/**
 * 函数方式增加下载器中间件 按添加顺序执行
 */
func (self *Spider) AddDlMiddlewareFunc(reqFunc func(req *core.Request) interface{}, respFunc func(resp *core.Response) interface{}) *Spider {
	dlm := middleware.NewDlMiddleware(reqFunc, respFunc)
	return self.AddDlMiddleware(dlm)
}

func (self *Spider) ReqManager(rm reqmnger.Ireqmnger) *Spider {
	self.SchedulerConfig.ReqManager(rm)
	return self
}

/**
 * 创建爬虫实例
 */
func NewSpider(name string) *Spider {
	sp := &Spider {
		name 			: name,
		stop			: make(chan bool),
	}

	sp.SchedulerConfig = &SchedulerConfig{}
	sp.DownloaderConfig = &DownloaderConfig{}

	sp.DlThreadNum(runtime.NumCPU())
	sp.ParseThreadNum(runtime.NumCPU())
	sp.DlDuration(time.Second * 2)

	return sp
}