package tezos

import (
	"context"
	"time"

	"github.com/guillaumedebavelaere/tezos-delegation/pkg/http"
)

// API describes the tezos API interface.
type API interface {
	http.Client
	ListDelegations(ctx context.Context, fromTimestamp *time.Time) ([]*Delegation, error)
}
