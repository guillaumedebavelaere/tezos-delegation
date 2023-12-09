package cron

import (
	"context"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/datastore"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos"
	"go.uber.org/zap"
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

	// list delegations which will be stored
	zap.L().Info("list delegations from tezos service...")
	delegations, err := c.tezosService.ListDelegations(ctx, nil)
	if err != nil {
		return err
	}

	zap.L().Info("store delegations in datastore...")
	err = c.datastore.StoreDelegations(ctx, delegations)
	if err != nil {
		return err
	}

	return nil
}
