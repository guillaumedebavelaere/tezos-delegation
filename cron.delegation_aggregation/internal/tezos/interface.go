package tezos

import (
	"context"
	"time"
)

type API interface {
	ListDelegations(ctx context.Context, fromTimestamp *time.Time) ([]*Delegation, error)
}
