package http

import "github.com/imroc/req/v3"

type Client interface {
	C() *req.Client
	Init() error
}
