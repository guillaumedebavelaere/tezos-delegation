package tezos

import (
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/http"
)

// Config defines tezos client configuration.
type Config struct {
	HTTP http.ClientConfig `mapstructure:",squash"`
}

// Client represents tezos client.
type Client struct {
	http.Client
	cfg *Config
}

// NewClient creates a new tezos client.
func NewClient(cfg *Config, options ...http.Option) API {
	return &Client{
		Client: http.NewClient(&cfg.HTTP, options...),
		cfg:    cfg,
	}
}

// Init initializes tezos client.
func (c *Client) Init() {
	c.Client.Init()
}
