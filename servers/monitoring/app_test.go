package monitoring_test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring"
	"github.com/patharetush/incident-pilot/servers/monitoring/config"
	"github.com/patharetush/incident-pilot/servers/monitoring/tools"
	"github.com/patharetush/incident-pilot/shared/mcptest"
)

func TestMonitoringMCPEndToEnd(t *testing.T) {
	app := monitoring.New(config.Default(), nil)
	session := mcptest.ConnectHTTP(t, app.MCPServer())
	ctx := t.Context()

	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(toolsResult.Tools) != len(monitoring.ToolNames()) {
		t.Fatalf("got %d tools, want %d", len(toolsResult.Tools), len(monitoring.ToolNames()))
	}

	t.Run("list_services", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tools.ListServices})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.ListServicesOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Count != 3 {
			t.Fatalf("got %d services, want 3", out.Count)
		}
	})

	t.Run("get_service_metrics", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.GetServiceMetrics,
			Arguments: map[string]any{"service": "payment-api"},
		})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.GetServiceMetricsOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Service != "payment-api" || len(out.Metrics) == 0 {
			t.Fatalf("unexpected metrics output: %+v", out)
		}
	})

	t.Run("list_alerts", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.ListAlerts,
			Arguments: map[string]any{"severity": "critical"},
		})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.ListAlertsOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Count != 2 {
			t.Fatalf("got %d critical alerts, want 2", out.Count)
		}
	})

	t.Run("get_alert", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.GetAlert,
			Arguments: map[string]any{"alert_id": "alert-001"},
		})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.GetAlertOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Alert.ID != "alert-001" {
			t.Fatalf("got alert %q", out.Alert.ID)
		}
	})

	t.Run("get_service_metrics_unknown_service", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.GetServiceMetrics,
			Arguments: map[string]any{"service": "does-not-exist"},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected tool error for unknown service")
		}
	})

	t.Run("service_catalog_resource", func(t *testing.T) {
		result, err := session.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "monitoring://catalog/services",
		})
		if err != nil || len(result.Contents) == 0 {
			t.Fatalf("ReadResource: err=%v contents=%d", err, len(result.Contents))
		}
	})

	t.Run("investigate_incident_prompt", func(t *testing.T) {
		result, err := session.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "investigate_incident",
			Arguments: map[string]string{"service": "payment-api"},
		})
		if err != nil || len(result.Messages) == 0 {
			t.Fatalf("GetPrompt: err=%v", err)
		}
	})
}
