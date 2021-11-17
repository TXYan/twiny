package middleware

import (
	"twiny/core"
)

// downloaderMiddleware 返回值 暂且借鉴scrapy的用法
type IDlMiddleware interface {
	/**
	 * 返回 nil 继续下一个中间件
	 * 返回 *core.Request 终止中间件调用 且将返回的request放入到队列中
	 * 返回 *core.Response 终止继续后续的中间件调用 已经激活的按倒序调用
	 */
	ProcessRequest(req *core.Request) interface{}
	/**
	 * 返回 nil 继续按参数resp调用下一个中间件
	 * 返回 *core.Request 终止中间件调用 且将返回的request放入到队列中
	 * 返回 *core.Response 按照返回的response继续调用下一个中间件
	 */
	ProcessResponse(resp *core.Response) interface{}
}

// 定义downloaderMiddleware 可以节省使用者定义新的结构体 类似Abstract
type DlMiddleware struct {
	processRequestFunc func(req *core.Request) interface{}
	processResponseFunc func(resp *core.Response) interface{}
}

func (self *DlMiddleware) ProcessRequest(req *core.Request) interface{} {
	if self.processRequestFunc != nil {
		return self.processRequestFunc(req)
	}
	return nil
}

func (self *DlMiddleware) ProcessResponse(resp *core.Response) interface{} {
	if self.processResponseFunc != nil {
		return self.processResponseFunc(resp)
	}
	return nil
}

func NewDlMiddleware(reqFunc func(req *core.Request) interface{}, respFunc func(resp *core.Response) interface{}) *DlMiddleware {
	return &DlMiddleware {
		processRequestFunc 	: reqFunc,
		processResponseFunc : respFunc,
	}
}