package downloader

import (
	"io/ioutil"
	"net/http"
	"twiny/core"
	"twiny/config"
	"twiny/middleware"
)

/**
 * 下载器
 */

type Downloader struct {
	client *http.Client
	dlConfig config.DownloaderConfig
}

func (self *Downloader) middlewareRequest(req *core.Request, dlms []middleware.IDlMiddleware) (int, *core.Request, *core.Response) {
	if len(dlms) <= 0 {
		return -1, nil, nil
	}


	for idx, middleware := range dlms {
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

	return len(dlms) - 1, nil, nil
}


/**
 * 下载器中间件按执行过ProcessRequest的中间件 倒序执行 ProcessResponse
 */
func (self *Downloader) middlwareResponse(resp *core.Response, dlms []middleware.IDlMiddleware, maxIdx int) (*core.Request, *core.Response) {

	// middlewares
	if len(dlms) <= 0 {
		return nil, resp
	}

	/**
	 * 返回 nil 继续按参数resp调用下一个中间件
	 * 返回 *core.Request 终止中间件调用 且将返回的request放入到队列中
	 * 返回 *core.Response 按照返回的response继续调用下一个中间件
	 */
	for idx := maxIdx; idx >= 0; idx-- {
		result := dlms[idx].ProcessResponse(resp)
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


func (self *Downloader) Download(req *core.Request) (*core.Request, *core.Response) {
	// middlewares
	dlms := self.dlConfig.GetDlMiddlewares()
	dlmIdx, newReq, resp := self.middlewareRequest(req, dlms)

	if newReq != nil {
		// 重新入队列
	 	return newReq, nil
	}

	// 如果已经有结果则不调用downloader
	if resp == nil {
		resp = download(self.client, req)
	}

	return self.middlwareResponse(resp, dlms, dlmIdx)

}

func download(client *http.Client, req *core.Request) *core.Response {
	resp, err := client.Do(req.HttpRequest)

	var bodyBytes []byte
	if err == nil {
		defer resp.Body.Close()
		
		bodyBytes, err = ioutil.ReadAll(resp.Body)
	}

	return &core.Response {
		HttpResponse: resp,
		Body		: bodyBytes,
		Error		: err,
		Request		: req,
	}
}

func NewDownloader(dlConfig config.DownloaderConfig) *Downloader {
	client := dlConfig.GetDlClient()
	if client == nil {
		client = http.DefaultClient
	}

	downloader := &Downloader {
		dlConfig: dlConfig,
		client	: client,
	}

	return downloader
}
