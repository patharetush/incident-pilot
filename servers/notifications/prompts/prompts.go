package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/notifications/service"
)

const DraftIncidentUpdate = "draft_incident_update"

func Register(server *mcp.Server, svc *service.Service) {
	server.AddPrompt(&mcp.Prompt{
		Name: DraftIncidentUpdate,
		Description: "Draft a stakeholder incident update with approval-aware channel guidance",
		Arguments: []*mcp.PromptArgument{
			{Name: "service", Description: "Affected service (e.g. payment-api)", Required: true},
			{Name: "severity", Description: "Incident severity (e.g. SEV-2)", Required: false},
		},
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		serviceName := strings.TrimSpace(req.Params.Arguments["service"])
		severity := strings.TrimSpace(req.Params.Arguments["severity"])
		if serviceName == "" {
			return nil, fmt.Errorf("service argument is required")
		}
		if severity == "" {
			severity = "SEV-2"
		}

		channels, err := svc.ListChannels(ctx)
		if err != nil {
			return nil, err
		}
		recent, err := svc.ListNotifications(ctx, "", serviceName, 5)
		if err != nil {
			return nil, err
		}
		pending, err := svc.ListPendingApprovals(ctx)
		if err != nil {
			return nil, err
		}

		var channelLines []string
		for _, ch := range channels {
			channelLines = append(channelLines, fmt.Sprintf("- [%s] %s (%s): %s", ch.ID, ch.Name, ch.Type, ch.Audience))
		}

		var recentLines []string
		for _, n := range recent {
			recentLines = append(recentLines, fmt.Sprintf("- [%s] %s → %s", n.Status, n.Subject, n.ChannelID))
		}
		if len(recentLines) == 0 {
			recentLines = append(recentLines, "- No notifications sent yet for this service")
		}

		text := fmt.Sprintf(`Draft an incident update for %q (%s).

Available channels:
%s

Recent notifications:
%s

Pending approvals: %d

Drafting guidelines:
1. Lead with customer impact and current status.
2. Use #incidents (ch-slack-incidents) for operational updates.
3. Route executive summaries (ch-email-leads) through approval workflow.
4. Queue via queue_notification with requires_approval=true for sensitive channels.
5. Never send PagerDuty pages without explicit operator approval.`,
			serviceName, severity,
			strings.Join(channelLines, "\n"),
			strings.Join(recentLines, "\n"),
			len(pending))

		return &mcp.GetPromptResult{
			Description: fmt.Sprintf("Incident update draft for %s", serviceName),
			Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
		}, nil
	})
}
