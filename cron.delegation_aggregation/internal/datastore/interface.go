package datastore

import (
	"context"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/datastore/model"
)

type Datastorer interface {
	StoreDelegations(ctx context.Context, delegations []*model.Delegation) error
	GetLatestDelegation(ctx context.Context) (*model.Delegation, error)
}
