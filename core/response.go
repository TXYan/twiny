package core

import (
	"net/http"
)

type HttpResponse = http.Response

type Response struct {
	*HttpResponse
	Body []byte    // 覆盖http.Response的Body
	Error error

	Request *Request
}

func NewResponse() *Response {
	return &Response {
		HttpResponse	: &HttpResponse{},
	}
}