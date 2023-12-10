package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config describes the mongo client configuration.
type Config struct {
	URI               string        `validate:"required"`
	ConnectTimeout    time.Duration `validate:"required"`
	HeartbeatInterval time.Duration `validate:"required"`
	Timeout           time.Duration `validate:"required"`
	Username          string        `validate:"required"`
	Password          string        `validate:"required"`
}

type client struct {
	cfg *Config
	c   *mongo.Client
}

// New creates a new Mongo database client.
func New(cfg *Config) Client {
	return &client{cfg: cfg}
}

// Init initialize the mongo client.
func (c *client) Init() error {
	mongoClient, err := mongo.Connect(
		context.Background(),
		options.Client().
			ApplyURI(c.cfg.URI).
			SetConnectTimeout(c.cfg.ConnectTimeout).
			SetHeartbeatInterval(c.cfg.HeartbeatInterval).
			SetTimeout(c.cfg.Timeout).
			SetAuth(options.Credential{
				AuthMechanism: "SCRAM-SHA-256",
				Username:      c.cfg.Username,
				Password:      c.cfg.Password,
			}),
	)
	if err != nil {
		return err
	}

	c.c = mongoClient

	return c.Ping()
}

// Ping ping mongo database.
func (c *client) Ping() error {
	return c.c.Ping(context.Background(), readpref.Primary())
}

// Close disconnect from mongo.
func (c *client) Close() error {
	return c.c.Disconnect(context.Background())
}

// C returns the mongo client.
func (c *client) C() *mongo.Client {
	return c.c
}
