package knowledge_test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/knowledge"
	"github.com/patharetush/incident-pilot/servers/knowledge/config"
	"github.com/patharetush/incident-pilot/servers/knowledge/tools"
	"github.com/patharetush/incident-pilot/shared/mcptest"
)

func TestKnowledgeMCPEndToEnd(t *testing.T) {
	app := knowledge.New(config.Default(), nil)
	session := mcptest.ConnectHTTP(t, app.MCPServer())
	ctx := t.Context()

	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(toolsResult.Tools) != len(knowledge.ToolNames()) {
		t.Fatalf("got %d tools, want %d", len(toolsResult.Tools), len(knowledge.ToolNames()))
	}

	t.Run("search_runbooks", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.SearchRunbooks,
			Arguments: map[string]any{"service": "payment-api", "query": "error"},
		})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.SearchRunbooksOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Count == 0 {
			t.Fatal("expected runbooks")
		}
	})

	t.Run("get_runbook", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.GetRunbook,
			Arguments: map[string]any{"runbook_id": "rb-001"},
		})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.GetRunbookOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if len(out.Runbook.Steps) == 0 {
			t.Fatal("expected runbook steps")
		}
	})

	t.Run("recommend_mitigation_prompt", func(t *testing.T) {
		result, err := session.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "recommend_mitigation",
			Arguments: map[string]string{"service": "payment-api", "symptom": "timeout"},
		})
		if err != nil || len(result.Messages) == 0 {
			t.Fatalf("GetPrompt: err=%v", err)
		}
	})
}
