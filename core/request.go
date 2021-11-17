package core

import (
	"bytes"
	"net/http"
)

type HttpRequest = http.Request

type Request struct {
	*HttpRequest
	RetryTimes		int
}

func Get(url string) *Request {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	return &Request {
		HttpRequest	: req,
	}
}

func PostString(url string, body string) *Request {
	return Post(url, []byte(body))
}

func Post(url string, body []byte) *Request {

	var reader *bytes.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil
	}

	return &Request {
		HttpRequest		: req,
	}
}