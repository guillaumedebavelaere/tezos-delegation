package mongo

import "go.mongodb.org/mongo-driver/mongo"
import mongosvc "github.com/guillaumedebavelaere/tezos-delegation/pkg/mongo"

const (
	database              = "tezos_delegation"
	collectionDelegations = "delegations"
)

// Datastore represents the implementation of the datastore with mongo.
type Datastore struct {
	client      mongosvc.Client
	delegations *mongo.Collection
}

// New create a new mongo datastore.
func New(client mongosvc.Client) *Datastore {
	return &Datastore{
		client: client,
	}
}

// Init initialize mongo datastore.
func (d *Datastore) Init() error {
	err := d.client.Init()
	if err != nil {

		return err
	}

	d.delegations = d.client.C().Database(database).Collection(collectionDelegations)

	return nil
}
