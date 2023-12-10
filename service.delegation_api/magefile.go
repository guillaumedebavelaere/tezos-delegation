//go:build mage

package main

import (
	//mage:import
	"github.com/guillaumedebavelaere/tezos-delegation/tools/mage/service"
)

func init() {
	service.Name = "delegation_api"
}
