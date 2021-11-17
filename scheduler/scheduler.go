package scheduler

import (
	"twiny/core"
)

type IScheduler interface {
	AddRequest(req *core.Request)
	FireRequest() <- chan *core.Request
	Close()
}