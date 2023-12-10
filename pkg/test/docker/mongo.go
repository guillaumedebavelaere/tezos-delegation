package docker

import (
	"fmt"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/guillaumedebavelaere/tezos-delegation/pkg/mongo"
)

const (
	// MongoUsername is the default username for the mongo container.
	MongoUsername = "admin"
	// MongoPassword is the default password for the mongo container.
	MongoPassword = "admin"

	mongoRepository = "mongo"
	mongoVersion    = "6.0.8"
)

// NewMongoContainer returns a new mongo container.
func NewMongoContainer() *Container {
	return &Container{
		Options: &dockertest.RunOptions{
			Repository: mongoRepository,
			Tag:        mongoVersion,
			Env: []string{
				"MONGO_INITDB_ROOT_USERNAME=" + MongoUsername,
				"MONGO_INITDB_ROOT_PASSWORD=" + MongoPassword,
			},
		},
		HostConfigOptions: []func(*docker.HostConfig){
			func(config *docker.HostConfig) {
				config.AutoRemove = true
				config.RestartPolicy = docker.RestartPolicy{Name: "no"}
			},
		},
		Retry: func(resource *dockertest.Resource) func() error {
			return func() error {
				mongoClient := mongo.New(&mongo.Config{
					URI:               EndpointMongo(resource),
					ConnectTimeout:    30 * time.Second,
					HeartbeatInterval: 30 * time.Second,
					Timeout:           30 * time.Second,
					Username:          MongoUsername,
					Password:          MongoPassword,
				})

				if err := mongoClient.Init(); err != nil {
					return err
				}

				return mongoClient.Ping()
			}
		},
		ExpireIn: 300,
	}
}

// EndpointMongo returns the endpoint of the mongo container.
func EndpointMongo(resource *dockertest.Resource) string {
	return fmt.Sprintf("mongodb://%s:%s@localhost:%s", MongoUsername, MongoPassword, resource.GetPort("27017/tcp"))
}
