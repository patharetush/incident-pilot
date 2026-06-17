package discovery

import (
	"context"
	"fmt"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/shared/logging"
)

// ServerCatalog describes a connected MCP server and its capabilities.
type ServerCatalog struct {
	Name  string       `json:"name"`
	URL   string       `json:"url"`
	Tools []*mcp.Tool  `json:"tools"`
}

// Pool manages MCP client sessions to downstream servers.
type Pool struct {
	mu       sync.RWMutex
	sessions map[string]*mcp.ClientSession
	catalog  map[string]*ServerCatalog
}

func NewPool() *Pool {
	return &Pool{
		sessions: make(map[string]*mcp.ClientSession),
		catalog:  make(map[string]*ServerCatalog),
	}
}

// Connect establishes sessions to all configured MCP server endpoints.
func (p *Pool) Connect(ctx context.Context, endpoints map[string]string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "incident-pilot-orchestrator",
		Version: "0.1.0",
	}, nil)

	var firstErr error
	for name, url := range endpoints {
		if url == "" {
			continue
		}
		session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: url}, nil)
		if err != nil {
			logging.L().Error().Err(err).Str("server", name).Str("url", url).Msg("MCP connect failed")
			if firstErr == nil {
				firstErr = fmt.Errorf("connect %s: %w", name, err)
			}
			continue
		}
		toolsResult, err := session.ListTools(ctx, nil)
		if err != nil {
			session.Close()
			if firstErr == nil {
				firstErr = fmt.Errorf("list tools %s: %w", name, err)
			}
			continue
		}
		p.sessions[name] = session
		p.catalog[name] = &ServerCatalog{Name: name, URL: url, Tools: toolsResult.Tools}
		logging.L().Info().
			Str("server", name).
			Int("tools", len(toolsResult.Tools)).
			Msg("MCP server discovered")
	}
	if len(p.sessions) == 0 {
		if firstErr != nil {
			return firstErr
		}
		return fmt.Errorf("no MCP servers connected")
	}
	return nil
}

// Catalog returns discovered server capabilities.
func (p *Pool) Catalog() []*ServerCatalog {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]*ServerCatalog, 0, len(p.catalog))
	for _, c := range p.catalog {
		out = append(out, c)
	}
	return out
}

// CallTool invokes a tool on a named MCP server.
func (p *Pool) CallTool(ctx context.Context, server, tool string, args map[string]any) (*mcp.CallToolResult, error) {
	p.mu.RLock()
	session, ok := p.sessions[server]
	p.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("server %q not connected", server)
	}
	return session.CallTool(ctx, &mcp.CallToolParams{Name: tool, Arguments: args})
}

// Close closes all MCP sessions.
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for name, session := range p.sessions {
		if err := session.Close(); err != nil {
			logging.L().Error().Err(err).Str("server", name).Msg("session close failed")
		}
	}
	p.sessions = make(map[string]*mcp.ClientSession)
	p.catalog = make(map[string]*ServerCatalog)
}

// HasServer reports whether a server is connected.
func (p *Pool) HasServer(name string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.sessions[name]
	return ok
}
