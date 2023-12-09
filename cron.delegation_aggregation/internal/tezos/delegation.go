package tezos

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

const (
	delegationsResource = "operations/delegations"
)

// Delegation represents the tezos delegation model.
type Delegation struct {
	Timestamp time.Time
	Amount    int64 `json:"amount"`
	Sender    struct {
		Address string `json:"address"`
	} `json:"sender"`
	Block string `json:"block"`
}

// ListDelegations returns delegations list.
func (c *Client) ListDelegations(ctx context.Context, fromTimestamp *time.Time) ([]*Delegation, error) {

	params := map[string]string{}
	// select only needed fields
	params["select"] = "timestamp,amount,sender,block"
	params["limit"] = "100"

	if fromTimestamp != nil {
		// filter by timestamp greater than fromTimestamp
		params["timestamp.gt"] = fromTimestamp.String()
		// couldn't sort by timestamp, but sort by id descending seems to be correlated, to confirm with tezos team.
		params["sort.asc"] = "id"

	} else {
		params["sort.desc"] = "id"
	}

	delegations := []*Delegation{}
	resp, err := c.client.R().
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

		return nil, fmt.Errorf("couldn't list delegations from tezos api: %s", resp.String())
	}

	return delegations, nil

}
