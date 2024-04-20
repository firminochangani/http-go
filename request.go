package http

import (
	"context"
	"net/url"
)

type Request struct {
	Method  string
	Headers Header
	URL     *url.URL
	Context context.Context
}