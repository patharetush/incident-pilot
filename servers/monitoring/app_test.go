package monitoring_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring"
	"github.com/patharetush/incident-pilot/servers/monitoring/config"
	"github.com/patharetush/incident-pilot/servers/monitoring/tools"
)

func TestMonitoringMCPEndToEnd(t *testing.T) {
	ctx := context.Background()
	cfg := config.Default()
	app := monitoring.New(cfg, nil)
	server := app.MCPServer()

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{JSONResponse: true})

	httpServer := httptest.NewServer(handler)
	t.Cleanup(httpServer.Close)

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "incident-pilot-test-client",
		Version: "0.1.0",
	}, nil)

	t.Log("httpServer.URL", httpServer.URL)

	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint: httpServer.URL,
	}, nil)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	t.Cleanup(func() { session.Close() })

	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(toolsResult.Tools) != len(monitoring.ToolNames()) {
		t.Fatalf("got %d tools, want %d", len(toolsResult.Tools), len(monitoring.ToolNames()))
	}

	t.Run("list_services", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.ListServices,
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if result.IsError {
			t.Fatalf("tool error: %v", result.Content)
		}

		var out tools.ListServicesOutput
		decodeStructuredOutput(t, result, &out)
		if out.Count != 3 {
			t.Fatalf("got %d services, want 3", out.Count)
		}
	})

	t.Run("get_service_metrics", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.GetServiceMetrics,
			Arguments: map[string]any{
				"service": "payment-api",
			},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if result.IsError {
			t.Fatalf("tool error: %v", result.Content)
		}

		var out tools.GetServiceMetricsOutput
		decodeStructuredOutput(t, result, &out)
		if out.Service != "payment-api" {
			t.Fatalf("got service %q, want payment-api", out.Service)
		}
		if len(out.Metrics) == 0 {
			t.Fatal("expected metrics for payment-api")
		}
	})

	t.Run("list_alerts", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.ListAlerts,
			Arguments: map[string]any{
				"severity": "critical",
			},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if result.IsError {
			t.Fatalf("tool error: %v", result.Content)
		}

		var out tools.ListAlertsOutput
		decodeStructuredOutput(t, result, &out)
		if out.Count != 2 {
			t.Fatalf("got %d critical alerts, want 2", out.Count)
		}
	})

	t.Run("get_alert", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.GetAlert,
			Arguments: map[string]any{
				"alert_id": "alert-001",
			},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if result.IsError {
			t.Fatalf("tool error: %v", result.Content)
		}

		var out tools.GetAlertOutput
		decodeStructuredOutput(t, result, &out)
		if out.Alert.ID != "alert-001" {
			t.Fatalf("got alert %q, want alert-001", out.Alert.ID)
		}
		if out.Alert.Service != "payment-api" {
			t.Fatalf("got service %q, want payment-api", out.Alert.Service)
		}
	})

	t.Run("get_service_metrics_unknown_service", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.GetServiceMetrics,
			Arguments: map[string]any{
				"service": "does-not-exist",
			},
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
		if err != nil {
			t.Fatalf("ReadResource: %v", err)
		}
		if len(result.Contents) == 0 {
			t.Fatal("expected resource contents")
		}
	})

	t.Run("investigate_incident_prompt", func(t *testing.T) {
		result, err := session.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "investigate_incident",
			Arguments: map[string]string{"service": "payment-api"},
		})
		if err != nil {
			t.Fatalf("GetPrompt: %v", err)
		}
		if len(result.Messages) == 0 {
			t.Fatal("expected prompt messages")
		}
	})
}

func decodeStructuredOutput(t *testing.T, result *mcp.CallToolResult, dest any) {
	t.Helper()
	if result.StructuredContent == nil {
		t.Fatal("expected structured content in tool result")
	}
	data, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		t.Fatalf("unmarshal structured content: %v", err)
	}
}
