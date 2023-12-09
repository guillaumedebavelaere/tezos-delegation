package datastore

import (
	"context"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos"
)

type Datastorer interface {
	StoreDelegations(ctx context.Context, delegations []*tezos.Delegation) error
}
