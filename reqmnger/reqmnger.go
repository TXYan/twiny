package reqmnger

import (
	"github.com/TXYan/twiny/core"
)

type Ireqmnger interface {
	/**
	 * 添加reqs的请求
	 */
	AddReqs(reqs ...*core.Request)

	/**
	 * 下一个需要处理的req
	 */
	NextReq() *core.Request
}