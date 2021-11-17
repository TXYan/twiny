package main

import (
	"fmt"
	"time"
	"strconv"
	"twiny/core"
	"twiny/spider"
)

func main() {

}


func demo() {
	spider.NewSpider("BaiduSpider", []*core.Request{core.Get("https://studygolang.com/articles/8865"),}).
		AddDlMiddlewareFunc(func(req *core.Request) interface{} {
			fmt.Println("DlMiddleware1 ProcessRequest")
			req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
			return false
		}, func(resp *core.Response) interface{} {
			fmt.Println("DlMiddleware1 ProcessResponse")
			return true
		}).
		AddDlMiddlewareFunc(func(req *core.Request) interface{} {
			fmt.Println("DlMiddleware2 ProcessRequest")
			
			if req.RetryTimes == 0 {
				req.RetryTimes++
				return req
			}

			resp := &core.Response{}
			resp.Status = "200 OK"
			resp.Request = req
			resp.Body = []byte("mock response")

			return resp

		}, func(resp *core.Response) interface{} {
			fmt.Println("DlMiddleware2 ProcessResponse")
			return true
		}).
		ParserFunc(func(resp *core.Response) []*core.Request {
			fmt.Println(resp.Request.URL.Path, ", response status:", resp.Status)
			return nil
		}).
		Crawl()
}

type baiduParser struct {
}


func (self *baiduParser) Parse(resp *core.Response) []*core.Request {
	fmt.Println(resp.Request.URL.Path, ", response status:", resp.Status)
	return nil
}



type baiduHeaderMiddleware struct {
}

func (self *baiduHeaderMiddleware) ProcessRequest(req *core.Request) {
	fmt.Println("baiduHeaderMiddleware ProcessRequest call")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
}

func (self *baiduHeaderMiddleware) ProcessResponse(resp *core.Response) {
	fmt.Println("baiduHeaderMiddleware ProcessResponse call")
}

