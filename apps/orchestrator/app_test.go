package orchestrator_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/apps/orchestrator"
	"github.com/patharetush/incident-pilot/apps/orchestrator/config"
	"github.com/patharetush/incident-pilot/apps/orchestrator/planner"
	"github.com/patharetush/incident-pilot/apps/orchestrator/sessions"
	"github.com/patharetush/incident-pilot/servers/deployments"
	deploycfg "github.com/patharetush/incident-pilot/servers/deployments/config"
	"github.com/patharetush/incident-pilot/servers/knowledge"
	knowcfg "github.com/patharetush/incident-pilot/servers/knowledge/config"
	"github.com/patharetush/incident-pilot/servers/logs"
	logcfg "github.com/patharetush/incident-pilot/servers/logs/config"
	"github.com/patharetush/incident-pilot/servers/monitoring"
	moncfg "github.com/patharetush/incident-pilot/servers/monitoring/config"
	"github.com/patharetush/incident-pilot/servers/notifications"
	notifcfg "github.com/patharetush/incident-pilot/servers/notifications/config"
)

func TestOrchestratorEndToEnd(t *testing.T) {
	endpoints := startMCPCluster(t)
	cfg := config.Default()
	cfg.MCP.Monitoring = endpoints[config.ServerMonitoring]
	cfg.MCP.Deployments = endpoints[config.ServerDeployments]
	cfg.MCP.Logs = endpoints[config.ServerLogs]
	cfg.MCP.Knowledge = endpoints[config.ServerKnowledge]
	cfg.MCP.Notifications = endpoints[config.ServerNotifications]

	app := orchestrator.New(cfg)
	t.Cleanup(app.Close)

	ctx := t.Context()
	if err := app.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	if len(app.Catalog()) != 5 {
		t.Fatalf("got %d servers, want 5", len(app.Catalog()))
	}

	session, err := app.RunInvestigation(ctx, "payment-api")
	if err != nil {
		t.Fatalf("RunInvestigation: %v", err)
	}

	if session.Status != sessions.StatusCompleted {
		t.Fatalf("got status %q, want completed", session.Status)
	}
	if session.Plan == nil || len(session.Plan.Steps) == 0 {
		t.Fatal("expected investigation plan")
	}
	if len(session.Evidence) == 0 {
		t.Fatal("expected collected evidence")
	}
	if len(session.Recommendations) == 0 {
		t.Fatal("expected recommendations")
	}

	completed := 0
	for _, step := range session.Plan.Steps {
		if step.Status == planner.StepCompleted {
			completed++
		}
	}
	if completed == 0 {
		t.Fatal("expected at least one completed plan step")
	}
}

func startMCPCluster(t *testing.T) map[string]string {
	t.Helper()

	servers := map[string]*mcp.Server{
		config.ServerMonitoring:    monitoring.New(moncfg.Default(), nil).MCPServer(),
		config.ServerDeployments:   deployments.New(deploycfg.Default(), nil).MCPServer(),
		config.ServerLogs:          logs.New(logcfg.Default(), nil).MCPServer(),
		config.ServerKnowledge:     knowledge.New(knowcfg.Default(), nil).MCPServer(),
		config.ServerNotifications: notifications.New(notifcfg.Default(), nil).MCPServer(),
	}

	endpoints := make(map[string]string, len(servers))
	for name, server := range servers {
		handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
			return server
		}, &mcp.StreamableHTTPOptions{JSONResponse: true})
		httpServer := httptest.NewServer(handler)
		t.Cleanup(httpServer.Close)
		endpoints[name] = httpServer.URL
	}
	return endpoints
}
