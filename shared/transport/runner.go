package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/shared/config"
	"github.com/patharetush/incident-pilot/shared/logging"
)

// Runner serves an MCP server over the configured transport.
type Runner struct {
	cfg    *config.Config
	server *mcp.Server
}

func NewRunner(cfg *config.Config, server *mcp.Server) *Runner {
	return &Runner{cfg: cfg, server: server}
}

func (r *Runner) Run(ctx context.Context) error {
	switch r.cfg.Transport.Mode {
	case config.TransportStdio:
		return r.runStdio(ctx)
	case config.TransportHTTP:
		return r.runHTTP(ctx)
	default:
		return fmt.Errorf("unsupported transport %q", r.cfg.Transport.Mode)
	}
}

func (r *Runner) runStdio(ctx context.Context) error {
	logging.L().Info().
		Str("server", r.cfg.Server.Name).
		Str("transport", config.TransportStdio).
		Msg("starting MCP server")
	return r.server.Run(ctx, &mcp.StdioTransport{})
}

func (r *Runner) runHTTP(ctx context.Context) error {
	mcpHandler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return r.server
	}, nil)

	var handler http.Handler = mcpHandler
	if r.cfg.Auth.Enabled {
		handler = wrapAuth(handler, r.cfg.Auth)
		logging.L().Info().Str("server", r.cfg.Server.Name).Msg("authorization middleware enabled")
	}

	httpServer := &http.Server{
		Addr:    r.cfg.Transport.HTTPAddr,
		Handler: handler,
	}

	logging.L().Info().
		Str("server", r.cfg.Server.Name).
		Str("transport", config.TransportHTTP).
		Str("addr", r.cfg.Transport.HTTPAddr).
		Msg("starting MCP server")

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), r.cfg.Transport.ShutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logging.L().Error().Err(err).Str("server", r.cfg.Server.Name).Msg("HTTP server shutdown failed")
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func wrapAuth(next http.Handler, authCfg config.AuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if authCfg.APIKey != "" {
			token := req.Header.Get("Authorization")
			expected := "Bearer " + authCfg.APIKey
			if token != expected {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, req)
	})
}
