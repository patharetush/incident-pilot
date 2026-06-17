package notifications

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/notifications/config"
	"github.com/patharetush/incident-pilot/servers/notifications/prompts"
	"github.com/patharetush/incident-pilot/servers/notifications/repository"
	"github.com/patharetush/incident-pilot/servers/notifications/repository/memory"
	"github.com/patharetush/incident-pilot/servers/notifications/resources"
	"github.com/patharetush/incident-pilot/servers/notifications/service"
	"github.com/patharetush/incident-pilot/servers/notifications/tools"
	"github.com/patharetush/incident-pilot/shared/transport"
)

type App struct {
	cfg    *config.Config
	server *mcp.Server
}

type Options struct {
	Repository repository.Repository
}

func New(cfg *config.Config, opts *Options) *App {
	var repo repository.Repository = memory.New()
	if opts != nil && opts.Repository != nil {
		repo = opts.Repository
	}

	svc := service.New(repo)
	server := mcp.NewServer(&mcp.Implementation{
		Name: cfg.Server.Name, Version: cfg.Server.Version,
	}, nil)

	tools.Register(server, svc)
	resources.Register(server, svc)
	prompts.Register(server, svc)

	return &App{cfg: cfg, server: server}
}

func (a *App) MCPServer() *mcp.Server { return a.server }

func (a *App) Run(ctx context.Context) error {
	return transport.NewRunner(a.cfg, a.server).Run(ctx)
}

func ToolNames() []string { return tools.Names() }
