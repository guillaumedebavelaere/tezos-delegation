package datastore

import (
	"context"

	"github.com/guillaumedebavelaere/tezos-delegation/pkg/tezos/datastore/model"
)

// Datastorer describes the datastore interface.
type Datastorer interface {
	StoreDelegations(ctx context.Context, delegations []*model.Delegation) error
	GetLatestDelegation(ctx context.Context) (*model.Delegation, error)
	GetDelegations(
		ctx context.Context,
		pageNumber, pageSize, year int,
	) ([]*model.Delegation, error)
	GetDelegationsCount(ctx context.Context, year int) (int, error)
}
