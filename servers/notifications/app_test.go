package notifications_test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/notifications"
	"github.com/patharetush/incident-pilot/servers/notifications/config"
	"github.com/patharetush/incident-pilot/servers/notifications/tools"
	"github.com/patharetush/incident-pilot/shared/mcptest"
)

func TestNotificationsMCPEndToEnd(t *testing.T) {
	app := notifications.New(config.Default(), nil)
	session := mcptest.ConnectHTTP(t, app.MCPServer())
	ctx := t.Context()

	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(toolsResult.Tools) != len(notifications.ToolNames()) {
		t.Fatalf("got %d tools, want %d", len(toolsResult.Tools), len(notifications.ToolNames()))
	}

	t.Run("list_channels", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tools.ListChannels})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.ListChannelsOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Count == 0 {
			t.Fatal("expected channels")
		}
	})

	t.Run("list_pending_approvals", func(t *testing.T) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tools.ListPendingApprovals})
		if err != nil || result.IsError {
			t.Fatalf("CallTool: err=%v", err)
		}
		var out tools.ListPendingApprovalsOutput
		mcptest.DecodeStructuredOutput(t, result, &out)
		if out.Count == 0 {
			t.Fatal("expected pending approvals")
		}
	})

	t.Run("queue_and_approve", func(t *testing.T) {
		queueResult, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.QueueNotification,
			Arguments: map[string]any{
				"channel_id":        "ch-slack-incidents",
				"subject":           "Test update",
				"body":              "Automated test notification",
				"service":           "payment-api",
				"requires_approval": true,
			},
		})
		if err != nil || queueResult.IsError {
			t.Fatalf("QueueNotification: err=%v", err)
		}
		var queued tools.QueueNotificationOutput
		mcptest.DecodeStructuredOutput(t, queueResult, &queued)

		approveResult, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: tools.ApproveNotification,
			Arguments: map[string]any{
				"notification_id": queued.Notification.ID,
				"approver":          "test-operator",
			},
		})
		if err != nil || approveResult.IsError {
			t.Fatalf("ApproveNotification: err=%v", err)
		}
		var approved tools.ApproveNotificationOutput
		mcptest.DecodeStructuredOutput(t, approveResult, &approved)
		if approved.Notification.Status != "sent" {
			t.Fatalf("expected sent status, got %s", approved.Notification.Status)
		}
	})

	t.Run("draft_incident_update_prompt", func(t *testing.T) {
		result, err := session.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "draft_incident_update",
			Arguments: map[string]string{"service": "payment-api", "severity": "SEV-2"},
		})
		if err != nil || len(result.Messages) == 0 {
			t.Fatalf("GetPrompt: err=%v", err)
		}
	})
}
