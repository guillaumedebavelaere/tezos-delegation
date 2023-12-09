package mongo

import (
	"context"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos"
)

func (d *Datastore) StoreDelegations(ctx context.Context, delegations []*tezos.Delegation) error {
	// Convert the Delegation slice to an interface slice for bulk insertion
	documents := make([]interface{}, len(delegations))
	for i, delegation := range delegations {
		documents[i] = delegation
	}

	// Perform bulk insertion
	_, err := d.delegations.InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	return nil
}
