package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/patharetush/incident-pilot/servers/logs"
	"github.com/patharetush/incident-pilot/servers/logs/config"
	"github.com/patharetush/incident-pilot/shared/logging"
)

func main() {
	cfg := config.Load()
	if err := logging.InitWithConfig(logging.Config{Filename: cfg.Log.Filename}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logging: %v\n", err)
		os.Exit(1)
	}
	defer logging.Close()

	app := logs.New(cfg, nil)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		logging.L().Fatal().Err(err).Msg("logs MCP server stopped")
	}
}
