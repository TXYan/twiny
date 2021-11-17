package downloader

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"twiny/core"
)

/**
 * 简单下载器实现
 */
type SimpleDownloader struct {
	client *http.Client
}


func (self *SimpleDownloader) Download(req *core.Request) *core.Response {
	fmt.Println("SimpleDownloader call")
	resp, err := self.client.Do(req.HttpRequest)

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


func NewSimpleDownloader(client *http.Client) *SimpleDownloader {

	if client == nil {
		client = http.DefaultClient
	}

	downloader := &SimpleDownloader {
		client	: client,
	}

	return downloader
}