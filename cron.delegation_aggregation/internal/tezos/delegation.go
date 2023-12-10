package tezos

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

const (
	delegationsResource = "operations/delegations"
)

// Sender describes the sender in tezos API.
type Sender struct {
	Address string `json:"address"`
}

// Delegation represents the tezos delegation model.
type Delegation struct {
	Timestamp time.Time
	Amount    int64  `json:"amount"`
	Sender    Sender `json:"sender"`
	Block     string `json:"block"`
}

// ListDelegations returns delegations list.
func (c *Client) ListDelegations(ctx context.Context, fromTimestamp *time.Time) ([]*Delegation, error) {
	params := map[string]string{}
	// select only needed fields
	params["select"] = "timestamp,amount,sender,block"
	params["limit"] = "100"

	if fromTimestamp != nil {
		// filter by timestamp greater than fromTimestamp
		params["timestamp.gt"] = fromTimestamp.UTC().Format(time.RFC3339)
		// couldn't sort by timestamp, but sort by id descending seems to be correlated, to confirm with tezos team.
		params["sort.asc"] = "id"
	} else {
		params["sort.desc"] = "id"
	}

	delegations := []*Delegation{}

	resp, err := c.C().R().
		SetContext(ctx).
		SetSuccessResult(&delegations).
		SetQueryParams(params).
		Get(delegationsResource)
	if err != nil {
		zap.L().Error("couldn't list delegations from tezos api", zap.Error(err))

		return nil, fmt.Errorf("couldn't list delegations from tezos api error: %w", err)
	}

	if resp.IsErrorState() {
		zap.L().Error(
			"couldn't list delegations from tezos api",
			zap.String("status", resp.GetStatus()),
			zap.String("body", resp.String()),
		)

		return nil, fmt.Errorf("couldn't list delegations from tezos api error: %s", resp.String())
	}

	return delegations, nil
}
