package docker

import (
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
)

// Container is a struct that contains all the information needed to run a docker container.
type Container struct {
	Options           *dockertest.RunOptions
	HostConfigOptions []func(*dc.HostConfig)
	Retry             func(resource *dockertest.Resource) func() error
	ExpireIn          uint
}
