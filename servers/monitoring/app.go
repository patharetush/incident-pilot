package monitoring

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring/config"
	"github.com/patharetush/incident-pilot/servers/monitoring/prompts"
	"github.com/patharetush/incident-pilot/servers/monitoring/repository"
	"github.com/patharetush/incident-pilot/servers/monitoring/repository/memory"
	"github.com/patharetush/incident-pilot/servers/monitoring/resources"
	"github.com/patharetush/incident-pilot/servers/monitoring/service"
	"github.com/patharetush/incident-pilot/servers/monitoring/tools"
	"github.com/patharetush/incident-pilot/shared/transport"
)

// App is the composition root for the monitoring MCP server.
type App struct {
	cfg    *config.Config
	server *mcp.Server
}

// Options customize App construction (e.g. inject a Prometheus-backed repository).
type Options struct {
	Repository repository.Repository
}

// New wires all layers and returns a ready-to-run application.
func New(cfg *config.Config, opts *Options) *App {
	var repo repository.Repository = memory.New()
	if opts != nil && opts.Repository != nil {
		repo = opts.Repository
	}

	svc := service.New(repo)
	server := mcp.NewServer(&mcp.Implementation{
		Name:    cfg.Server.Name,
		Version: cfg.Server.Version,
	}, nil)

	tools.Register(server, svc)
	resources.Register(server, svc)
	prompts.Register(server, svc)

	return &App{cfg: cfg, server: server}
}

// MCPServer exposes the underlying MCP server (for tests and embedding).
func (a *App) MCPServer() *mcp.Server {
	return a.server
}

// Run starts the configured transport.
func (a *App) Run(ctx context.Context) error {
	return transport.NewRunner(a.cfg, a.server).Run(ctx)
}

// ToolNames returns registered MCP tool names.
func ToolNames() []string {
	return tools.Names()
}
