package mongo

import (
	"context"
	"errors"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/datastore/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *Datastore) StoreDelegations(ctx context.Context, delegations []*model.Delegation) error {
	// Convert the Delegation slice to an interface slice for bulk insertion
	documents := make([]interface{}, len(delegations))
	for i, delegation := range delegations {
		documents[i] = delegation
	}

	// Perform bulk insertion
	_, err := d.delegations.UpdateMany()InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	return nil
}

// GetLatestDelegation get the latest delegation in database (with the more recent timestamp).
func (d *Datastore) GetLatestDelegation(ctx context.Context) (*model.Delegation, error) {
	// An empty filter matches all documents
	filter := bson.D{{}}
	// sort to find the document with the latest timestamp
	sort := options.FindOne().SetSort(bson.D{{"timestamp", -1}})

	var result *model.Delegation
	err := d.delegations.FindOne(ctx, filter, sort).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return result, nil
}
