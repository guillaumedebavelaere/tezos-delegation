package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/cron"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/datastore/mongo"
	"github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/tezos"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/config"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/log"
	mongosvc "github.com/guillaumedebavelaere/tezos-delegation/pkg/mongo"
)

const appName = "delegation_aggregation"

//nolint:funlen
func run() int {
	log.SetDefaultZap()

	var cfg struct {
		Debug bool
		API   struct {
			Tezos tezos.Config
		}
		Datastore struct {
			Mongo mongosvc.Config
		}
	}

	// parse yaml config
	if err := config.Parse(appName, &cfg); err != nil {
		zap.L().Error("couldn't parse config", zap.Error(err))

		return 1
	}

	if err := config.Validate(cfg); err != nil {
		zap.L().Error("invalid config", zap.Error(err))

		return 1
	}

	log.Configure(cfg.Debug)

	tezosService := tezos.NewClient(&cfg.API.Tezos)

	err := tezosService.Init()
	if err != nil {
		zap.L().Error("couldn't create new tezos service", zap.Error(err))

		return 1
	}

	mongoClient := mongosvc.New(&cfg.Datastore.Mongo)
	datastore := mongo.New(mongoClient)

	err = datastore.Init()
	if err != nil {
		zap.L().Error(
			"couldn't initialize datastore",
			zap.Error(err),
		)

		return 1
	}

	// Create new delegation aggregation cron
	c := cron.New(tezosService, datastore)

	// run cronjob
	err = c.Run()
	if err != nil {
		zap.L().Error(
			"couldn't run delegation aggregation cron",
			zap.Error(err),
		)

		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
