package cron

import (
	"context"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/datastore"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/datastore/model"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos"
	"go.uber.org/zap"
	"time"
)

// Cron describes the delegation aggregation Cron.
type Cron struct {
	tezosService tezos.API
	datastore    datastore.Datastorer
}

// New creates a new Cron.
func New(tezosService tezos.API, datastore datastore.Datastorer) *Cron {
	return &Cron{
		tezosService: tezosService,
		datastore:    datastore,
	}
}

// Run polls the delegations from tezos API and store them in datastore.
func (c *Cron) Run() error {
	ctx := context.Background()

	latestDelegation, err := c.datastore.GetLatestDelegation(ctx)
	if err != nil {
		zap.L().Error("couldn't get latest delegation from datastore", zap.Error(err))

		return err
	}

	zap.L().Info("list delegations from tezos service ...")
	// list delegations which will be stored
	var latestTimestamp *time.Time
	if latestDelegation != nil {
		latestTimestamp = &latestDelegation.Timestamp
		zap.L().Info("from timestamp", zap.Any("latestTimestamp", latestTimestamp))
	}

	delegations, err := c.tezosService.ListDelegations(ctx, latestTimestamp)
	if err != nil {
		return err
	}

	if len(delegations) == 0 {
		zap.L().Info("no new delegations found")
		return nil
	}

	zap.L().Info("found", zap.Int("delegations", len(delegations)))
	zap.L().Info("store delegations in datastore...")
	err = c.storeDelegations(ctx, delegations)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cron) storeDelegations(ctx context.Context, tezosDelegations []*tezos.Delegation) error {
	delegationModels := make([]*model.Delegation, len(tezosDelegations))
	for i, tezosDelegation := range tezosDelegations {
		delegationModels[i] = &model.Delegation{
			Delegator: tezosDelegation.Sender.Address,
			Block:     tezosDelegation.Block,
			Amount:    tezosDelegation.Amount,
			Timestamp: tezosDelegation.Timestamp,
		}
	}

	return c.datastore.StoreDelegations(ctx, delegationModels)
}
