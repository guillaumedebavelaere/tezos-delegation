package http

import (
	"net/http"
	"time"

	"github.com/imroc/req/v3"

	"github.com/guillaumedebavelaere/tezos-delegation/pkg/option"
)

// ClientConfig represents the configuration used when creating a new HTTP Client.
type ClientConfig struct {
	Debug   bool
	BaseURL string        `validate:"required,url"`
	Timeout time.Duration `validate:"required"`
}

// Option custom option type to handle none exported struct.
type Option option.Option[*client]

type client struct {
	cfg     *ClientConfig
	c       *req.Client
	options []Option
}

// NewClient creates a new HTTP client base service.
func NewClient(cfg *ClientConfig, options ...Option) Client {
	return &client{
		cfg:     cfg,
		c:       req.C(),
		options: options,
	}
}

// Init initializes the http client.
func (c *client) Init() {
	// Set options
	for _, o := range c.options {
		o(c)
	}

	c.c = c.c.SetBaseURL(c.cfg.BaseURL).
		SetTimeout(c.cfg.Timeout)

	if c.cfg.Debug {
		c.c = c.c.DevMode()
	}
}

// C returns http client.
func (c *client) C() *req.Client {
	return c.c
}

// WithTransport is a Client option to customize http client Transport.
func WithTransport(t http.RoundTripper) func(*client) {
	return func(c *client) {
		c.c.GetClient().Transport = t
	}
}

// WithUnmarshaller is a Client option to customize unmarshaller.
func WithUnmarshaller(unmarshaller func(data []byte, v interface{}) error) func(*client) {
	return func(c *client) { c.c.SetJsonUnmarshal(unmarshaller) }
}
