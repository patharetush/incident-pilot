package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/deployments/service"
)

const CorrelateDeployment = "correlate_deployment"

func Register(server *mcp.Server, svc *service.Service) {
	server.AddPrompt(&mcp.Prompt{
		Name:        CorrelateDeployment,
		Description: "Correlate recent deployments with an active incident on a service",
		Arguments: []*mcp.PromptArgument{{
			Name: "service", Description: "Affected service (e.g. payment-api)", Required: true,
		}},
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		serviceName := strings.TrimSpace(req.Params.Arguments["service"])
		if serviceName == "" {
			return nil, fmt.Errorf("service argument is required")
		}

		deployments, err := svc.ListDeployments(ctx, serviceName, 5)
		if err != nil {
			return nil, err
		}
		changes, err := svc.ListRecentChanges(ctx, serviceName, 120)
		if err != nil {
			return nil, err
		}

		var depLines []string
		for _, d := range deployments {
			depLines = append(depLines, fmt.Sprintf("- %s %s (%s) at %s: %s",
				d.ID, d.Version, d.Status, d.DeployedAt.Format("15:04 UTC"), d.Summary))
		}
		if len(depLines) == 0 {
			depLines = append(depLines, "- No recent deployments found")
		}

		var changeLines []string
		for _, c := range changes {
			changeLines = append(changeLines, fmt.Sprintf("- [%s] %s", c.Type, c.Description))
		}
		if len(changeLines) == 0 {
			changeLines = append(changeLines, "- No recent changes in the lookback window")
		}

		text := fmt.Sprintf(`Correlate deployments with the incident on %q.

Recent deployments:
%s

Recent changes (last 120 minutes):
%s

Correlation checklist:
1. Compare incident start time with latest deployment timestamps.
2. Inspect config/code/migration changes for risky diffs.
3. Check if rollback candidates exist and are safe.
4. Gather deployment evidence before recommending rollback.`,
			serviceName, strings.Join(depLines, "\n"), strings.Join(changeLines, "\n"))

		return &mcp.GetPromptResult{
			Description: fmt.Sprintf("Deployment correlation for %s", serviceName),
			Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
		}, nil
	})
}
