package http_go

import (
	"context"
	"net/url"
)

type Request struct {
	Method     string
	Headers    Header
	URL        *url.URL
	RequestURI string
	Proto      string
	Host       string

	ctx context.Context
}

func (r *Request) Context() context.Context {
	return r.ctx
}
