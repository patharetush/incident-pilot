package main

import (
	"context"
	"fmt"
	"os"

	"github.com/patharetush/incident-pilot/apps/orchestrator"
	"github.com/patharetush/incident-pilot/apps/orchestrator/config"
	"github.com/patharetush/incident-pilot/shared/logging"
)

func main() {
	cfg := config.Load()

	if err := logging.InitWithConfig(logging.Config{Filename: cfg.Log.Filename}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logging: %v\n", err)
		os.Exit(1)
	}
	defer logging.Close()

	app := orchestrator.New(cfg)
	defer app.Close()

	ctx := context.Background()

	logging.L().Info().Msg("connecting to MCP servers")
	if err := app.Connect(ctx); err != nil {
		logging.L().Fatal().Err(err).Msg("MCP discovery failed")
	}

	for _, catalog := range app.Catalog() {
		logging.L().Info().
			Str("server", catalog.Name).
			Int("tools", len(catalog.Tools)).
			Str("url", catalog.URL).
			Msg("discovered MCP server")
	}

	session, err := app.RunInvestigation(ctx, cfg.Orchestrator.DefaultService)
	if err != nil {
		logging.L().Fatal().Err(err).Msg("investigation failed")
	}

	orchestrator.PrintReport(session)
}
