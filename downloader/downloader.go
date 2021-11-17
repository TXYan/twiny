package downloader

import (
	"twiny/core"
)

type IDownloader interface {
	Download(req *core.Request) *core.Response
}