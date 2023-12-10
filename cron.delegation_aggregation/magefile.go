//go:build mage

package main

import (
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/mage/gen"
	//mage:import
	"github.com/guillaumedebavelaere/tezos-delegation/tools/mage/service"
)

func init() {
	service.Name = "delegation_aggregation"
	service.GenFiles = []*gen.File{
		{
			Name:      "tezos",
			Type:      gen.Mock,
			Dest:      "./internal/tezos",
			Interface: []string{"API"},
			Pkg:       "github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos",
		},
	}
}
