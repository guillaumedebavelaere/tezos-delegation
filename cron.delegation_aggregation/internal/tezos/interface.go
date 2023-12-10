package tezos

import (
	"context"
	"time"
)

// API describes the tezos API interface.
type API interface {
	ListDelegations(ctx context.Context, fromTimestamp *time.Time) ([]*Delegation, error)
}
