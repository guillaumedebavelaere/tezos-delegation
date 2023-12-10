package tezos

import (
	"time"

	"github.com/imroc/req/v3"
)

// ClientConfig represents the configuration used when creating a new HTTP Client.
type ClientConfig struct {
	Debug   bool
	BaseURL string        `validate:"required,url"`
	Timeout time.Duration `validate:"required"`
}

// Client defines tezos client.
type Client struct {
	client *req.Client
}

// New returns a new tezos client.
func New(cfg *ClientConfig) (*Client, error) {
	c := &Client{
		client: req.NewClient(),
	}

	c.client = c.client.SetBaseURL(cfg.BaseURL).
		SetTimeout(cfg.Timeout)

	if cfg.Debug {
		c.client = c.client.DevMode()
	}

	return c, nil
}
