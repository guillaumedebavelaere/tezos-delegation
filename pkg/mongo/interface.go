package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Client defines a standard interface for Mongo database client.
type Client interface {
	C() *mongo.Client
	Ping() error
	Init() error
	Close() error
}
