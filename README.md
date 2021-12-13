# twiny
go spider


DEMO


``` golang

/**
 * 可以函数式编写DlMiddleware(下载器中间件), Parser解析器
 * 可以自己实现middleware.IDlMiddleware, parser.IParser
 */

spider.NewSpider("CsdnSpider").
	StartReqs([]*core.Request{core.Get("https://blog.csdn.net/qq_33513250/article/details/102989256"),}).
	AddDlMiddlewareFunc(func(req *core.Request) interface{} {
		fmt.Println("DlMiddleware1 ProcessRequest")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
		return nil
	}, func(resp *core.Response) interface{} {
		fmt.Println("DlMiddleware1 ProcessResponse")
		return nil
	}).
	AddDlMiddlewareFunc(func(req *core.Request) interface{} {
		fmt.Println("DlMiddleware2 ProcessRequest")
		
		if req.RetryTimes == 0 {
			req.RetryTimes++
			return req
		}

		return nil
	}, func(resp *core.Response) interface{} {
		fmt.Println("DlMiddleware2 ProcessResponse")
		return nil
	}).
	ParserFunc(func(resp *core.Response) []*core.Request {
		fmt.Println(resp.Request.URL.Path, ", html:", string(resp.Body))
		return nil
	}).
	DlDuration(5 * time.Second).
	DlThreadNum(2).
	ParseThreadNum(2).
	Crawl()

```
