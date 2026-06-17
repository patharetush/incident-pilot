package knowledge

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/knowledge/config"
	"github.com/patharetush/incident-pilot/servers/knowledge/prompts"
	"github.com/patharetush/incident-pilot/servers/knowledge/repository"
	"github.com/patharetush/incident-pilot/servers/knowledge/repository/memory"
	"github.com/patharetush/incident-pilot/servers/knowledge/resources"
	"github.com/patharetush/incident-pilot/servers/knowledge/service"
	"github.com/patharetush/incident-pilot/servers/knowledge/tools"
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
