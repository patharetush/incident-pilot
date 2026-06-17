package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring/service"
)

const InvestigateIncident = "investigate_incident"

// Register attaches monitoring MCP prompts to the server.
func Register(server *mcp.Server, svc *service.Service) {
	server.AddPrompt(&mcp.Prompt{
		Name:        InvestigateIncident,
		Description: "Structured incident investigation workflow for a monitored service",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "service",
				Description: "Service to investigate (e.g. payment-api)",
				Required:    true,
			},
		},
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		serviceName := strings.TrimSpace(req.Params.Arguments["service"])
		if serviceName == "" {
			return nil, fmt.Errorf("service argument is required")
		}

		svcDetail, err := svc.GetService(ctx, serviceName)
		if err != nil {
			return nil, err
		}

		_, metrics, err := svc.GetServiceMetrics(ctx, serviceName)
		if err != nil {
			return nil, err
		}

		alerts, err := svc.ListAlerts(ctx, "")
		if err != nil {
			return nil, err
		}

		var related []string
		for _, alert := range alerts {
			if strings.EqualFold(alert.Service, serviceName) {
				related = append(related, fmt.Sprintf("- [%s] %s: %s", alert.Severity, alert.ID, alert.Title))
			}
		}
		if len(related) == 0 {
			related = append(related, "- No active alerts for this service")
		}

		var metricLines []string
		for _, m := range metrics {
			metricLines = append(metricLines, fmt.Sprintf("- %s: %.2f %s", m.Name, m.Value, m.Unit))
		}
		if len(metricLines) == 0 {
			metricLines = append(metricLines, "- No metrics available")
		}

		text := fmt.Sprintf(`Investigate the incident affecting service %q.

Current status: %s
Summary: %s
Environment: %s

Active alerts:
%s

Key metrics:
%s

Investigation steps:
1. Confirm blast radius using list_services and list_alerts tools.
2. Pull detailed metrics with get_service_metrics for correlated signals.
3. Cross-reference recent deployments and upstream dependencies.
4. Document evidence and recommend mitigations before any remediation.`,
			svcDetail.Name,
			svcDetail.Status,
			svcDetail.Summary,
			svcDetail.Environment,
			strings.Join(related, "\n"),
			strings.Join(metricLines, "\n"),
		)

		return &mcp.GetPromptResult{
			Description: fmt.Sprintf("Incident investigation plan for %s", serviceName),
			Messages: []*mcp.PromptMessage{
				{
					Role:    "user",
					Content: &mcp.TextContent{Text: text},
				},
			},
		}, nil
	})
}
