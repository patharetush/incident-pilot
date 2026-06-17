package logs_test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/logs"
	"github.com/patharetush/incident-pilot/servers/logs/config"
	"github.com/patharetush/incident-pilot/servers/logs/tools"
	"github.com/patharetush/incident-pilot/shared/mcptest"
)

func TestLogsMCPEndToEnd(t *testing.T) {
	app := logs.New(config.Default(), nil)
	session := mcptest.ConnectHTTP(t, app.MCPServer())
	ctx := t.Context()

	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(toolsResult.Tools) != len(logs.ToolNames()) {
		t.Fatalf("got %d tools, want %d", len(toolsResult.Tools), len(logs.ToolNames()))
	}

	t.Run("search_logs", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.SearchLogs,
			Arguments: map[string]any{"service": "payment-api", "level": "error"},
		})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.SearchLogsOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Count == 0 {
			t.Fatal("expected error logs")
		}
	})

	t.Run("get_log_context", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.GetLogContext,
			Arguments: map[string]any{"trace_id": "trace-abc123"},
		})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.GetLogContextOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Count < 2 {
			t.Fatalf("expected correlated logs, got %d", out.Count)
		}
	})

	t.Run("analyze_logs_prompt", func(t *testing.T) {
		result, err := session.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "analyze_logs", Arguments: map[string]string{"service": "payment-api"},
		})
		if err != nil || len(result.Messages) == 0 {
			t.Fatalf("GetPrompt: err=%v", err)
		}
	})
}
