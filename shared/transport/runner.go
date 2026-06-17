package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring/config"
	"github.com/patharetush/incident-pilot/shared/logging"
)

// Runner serves the MCP server over the configured transport.
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
	logging.L().Info().Str("transport", config.TransportStdio).Msg("starting monitoring MCP server")
	return r.server.Run(ctx, &mcp.StdioTransport{})
}

func (r *Runner) runHTTP(ctx context.Context) error {
	mcpHandler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return r.server
	}, nil)

	var handler http.Handler = mcpHandler
	if r.cfg.Auth.Enabled {
		handler = wrapAuth(handler, r.cfg.Auth)
		logging.L().Info().Msg("authorization middleware enabled")
	}

	httpServer := &http.Server{
		Addr:    r.cfg.Transport.HTTPAddr,
		Handler: handler,
	}

	logging.L().Info().
		Str("transport", config.TransportHTTP).
		Str("addr", r.cfg.Transport.HTTPAddr).
		Msg("starting monitoring MCP server")

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), r.cfg.Transport.ShutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logging.L().Error().Err(err).Msg("HTTP server shutdown failed")
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// wrapAuth is a placeholder for future bearer-token or OAuth middleware.
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
