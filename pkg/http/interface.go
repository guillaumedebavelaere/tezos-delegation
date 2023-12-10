package http

import "github.com/imroc/req/v3"

// Client interface for http client.
type Client interface {
	C() *req.Client
	Init() error
}
