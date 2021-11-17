package parser

import (
	"twiny/core"
)

type IParser interface {
	Parse(resp *core.Response) []*core.Request
}

// 定义parser 可以节省使用者定义结构体 类似Abstract
type Parser struct {
	parseFunc func(resp *core.Response) []*core.Request
}

func (self *Parser) Parse(resp *core.Response) []*core.Request {
	if self.parseFunc != nil {
		return self.parseFunc(resp)
	}
	return nil
}

func NewParser(parseFunc func(resp *core.Response) []*core.Request) *Parser {
	return &Parser {
		parseFunc 	: parseFunc,
	}
}