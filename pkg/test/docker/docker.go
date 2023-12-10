package docker

import (
	"github.com/hashicorp/go-multierror"
	"github.com/ory/dockertest/v3"
	"go.uber.org/zap"
)

const defaultExpireInSeconds = 120

// Docker is a struct that contains all the information needed to run a docker container.
type Docker struct {
	pool       *dockertest.Pool
	containers []*Container
	resources  []*dockertest.Resource
}

// New returns a new Docker struct.
func New(containers ...*Container) (*Docker, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		zap.L().Error("couldn't construct pool", zap.Error(err))

		return nil, err
	}

	if err = pool.Client.Ping(); err != nil {
		zap.L().Error("couldn't ping docker", zap.Error(err))

		return nil, err
	}

	return &Docker{
		pool: pool,
		// network:     network,
		containers: containers,
		// credentials: NewCredentials(),
		resources: []*dockertest.Resource{},
	}, nil
}

// Start starts all the containers.
func (d *Docker) Start() error {
	for _, container := range d.containers {
		if container.Options.Auth.Username == "" {
			zap.L().Info("use docker hub registry credentials", zap.String("username", container.Options.Auth.Username))
		} else {
			zap.L().Info("use other registry with credentials", zap.String("username", container.Options.Auth.Username))
		}

		resource, err := d.pool.RunWithOptions(container.Options, container.HostConfigOptions...)
		if err != nil {
			zap.L().Error("couldn't start resource", zap.String("name", container.Options.Name), zap.Error(err))

			return err
		}

		if container.Retry == nil {
			container.Retry = func(resource *dockertest.Resource) func() error {
				return func() error {
					return nil
				}
			}
		}

		if err = d.pool.Retry(container.Retry(resource)); err != nil {
			zap.L().Error("couldn't connect to docker retry failed", zap.String("name", container.Options.Name), zap.Error(err))

			return err
		}

		expireIn := uint(defaultExpireInSeconds)
		if container.ExpireIn > 0 {
			expireIn = container.ExpireIn
		}

		if err := resource.Expire(expireIn); err != nil {
			zap.L().Error("couldn't set expiry", zap.String("name", container.Options.Name), zap.Error(err))

			return err
		}

		d.resources = append(d.resources, resource)
	}

	return nil
}

// GetResources returns a resource by name.
func (d *Docker) GetResources() []*dockertest.Resource {
	return d.resources
}

// Stop stops all the containers and the network.
func (d *Docker) Stop() error {
	var me *multierror.Error

	for _, resource := range d.resources {
		if err := d.pool.Purge(resource); err != nil {
			zap.L().Error("couldn't purge resource", zap.Error(err))

			me = multierror.Append(me, err)
		}
	}

	return me.ErrorOrNil()
}
