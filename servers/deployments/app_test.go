package deployments_test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/deployments"
	"github.com/patharetush/incident-pilot/servers/deployments/config"
	"github.com/patharetush/incident-pilot/servers/deployments/tools"
	"github.com/patharetush/incident-pilot/shared/mcptest"
)

func TestDeploymentsMCPEndToEnd(t *testing.T) {
	app := deployments.New(config.Default(), nil)
	session := mcptest.ConnectHTTP(t, app.MCPServer())
	ctx := t.Context()

	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(toolsResult.Tools) != len(deployments.ToolNames()) {
		t.Fatalf("got %d tools, want %d", len(toolsResult.Tools), len(deployments.ToolNames()))
	}

	t.Run("list_deployments", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tools.ListDeployments})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v isError=%v", err, result != nil && result.IsError)
		}
		var out tools.ListDeploymentsOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Count == 0 {
			t.Fatal("expected deployments")
		}
	})

	t.Run("get_deployment", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.GetDeployment,
			Arguments: map[string]any{"deployment_id": "dep-001"},
		})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.GetDeploymentOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Deployment.Service != "payment-api" {
			t.Fatalf("got service %q", out.Deployment.Service)
		}
	})

	t.Run("correlate_deployment_prompt", func(t *testing.T) {
		result, err := session.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "correlate_deployment", Arguments: map[string]string{"service": "payment-api"},
		})
		if err != nil || len(result.Messages) == 0 {
			t.Fatalf("GetPrompt: err=%v messages=%d", err, len(result.Messages))
		}
	})
}
