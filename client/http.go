package client

import "github.com/valyala/fasthttp"

type HTTPClient interface {
	Do(*fasthttp.Request, *fasthttp.Response) error
}

type HTTPClientFunc func(*fasthttp.Request, *fasthttp.Response) error

func (f HTTPClientFunc) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	return f(req, resp)
}

var _ HTTPClient = DefaultHTTPClient

var DefaultHTTPClient = &fasthttp.Client{
	NoDefaultUserAgentHeader: true,
	RetryIfErr: func(request *fasthttp.Request, attempts int, err error) (resetTimeout bool, retry bool) {
		return false, false
	},
}
