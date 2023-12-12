package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/guillaumedebavelaere/tezos-delegation/pkg/config"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/log"
	mongosvc "github.com/guillaumedebavelaere/tezos-delegation/pkg/mongo"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/tezos/datastore/mongo"
	"github.com/guillaumedebavelaere/tezos-delegation/service.delegation_api/internal/delegation"
)

const appName = "delegation_api"

//nolint:funlen
func main() {
	log.SetDefaultZap()

	var cfg struct {
		Debug     bool
		Addr      string
		Datastore struct {
			Mongo mongosvc.Config
		}
	}

	// parse yaml config
	if err := config.Parse(appName, &cfg); err != nil {
		zap.L().Error("couldn't parse config", zap.Error(err))
		os.Exit(1)
	}

	if err := config.Validate(cfg); err != nil {
		zap.L().Error("invalid config", zap.Error(err))
		os.Exit(1)
	}

	log.Configure(cfg.Debug)

	mongoClient := mongosvc.New(&cfg.Datastore.Mongo)
	datastore := mongo.New(mongoClient)

	if err := datastore.Init(); err != nil {
		zap.L().Error(
			"couldn't initialize datastore",
			zap.Error(err),
		)

		os.Exit(1)
	}

	defer func(datastore *mongo.Datastore) {
		if err := datastore.Close(); err != nil {
			zap.L().Error(
				"couldn't close datastore",
				zap.Error(err),
			)
		}
	}(datastore)

	apiDelegationHandler := delegation.New(datastore)

	http.HandleFunc("/xtz/delegations", apiDelegationHandler.GetDelegationsHandler)

	zap.L().Info("server started and listening", zap.String("addr", cfg.Addr))

	// Create a new HTTP server with custom timeouts
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		zap.L().Error("server closed")
	} else if err != nil {
		zap.L().Error("error starting server", zap.Error(err))
		panic(err)
	}
}
