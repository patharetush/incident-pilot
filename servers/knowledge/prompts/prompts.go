package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/knowledge/service"
)

const RecommendMitigation = "recommend_mitigation"

func Register(server *mcp.Server, svc *service.Service) {
	server.AddPrompt(&mcp.Prompt{
		Name: RecommendMitigation,
		Description: "Recommend mitigations using runbooks and similar past incidents",
		Arguments: []*mcp.PromptArgument{{
			Name: "service", Description: "Affected service (e.g. payment-api)", Required: true,
		}, {
			Name: "symptom", Description: "Primary symptom (e.g. high error rate)", Required: false,
		}},
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		serviceName := strings.TrimSpace(req.Params.Arguments["service"])
		symptom := strings.TrimSpace(req.Params.Arguments["symptom"])
		if serviceName == "" {
			return nil, fmt.Errorf("service argument is required")
		}

		query := symptom
		if query == "" {
			query = "error"
		}

		runbooks, err := svc.SearchRunbooks(ctx, query, serviceName, 5)
		if err != nil {
			return nil, err
		}
		incidents, err := svc.SearchPastIncidents(ctx, query, serviceName, 3)
		if err != nil {
			return nil, err
		}

		var rbLines, incLines []string
		for _, rb := range runbooks {
			rbLines = append(rbLines, fmt.Sprintf("- [%s] %s: %s", rb.ID, rb.Title, rb.Summary))
		}
		if len(rbLines) == 0 {
			rbLines = append(rbLines, "- No matching runbooks; search broader knowledge base")
		}
		for _, inc := range incidents {
			incLines = append(incLines, fmt.Sprintf("- [%s] %s → %s", inc.ID, inc.Title, inc.Resolution))
		}
		if len(incLines) == 0 {
			incLines = append(incLines, "- No similar past incidents found")
		}

		text := fmt.Sprintf(`Recommend mitigations for %q (symptom: %q).

Matching runbooks:
%s

Similar past incidents:
%s

Mitigation workflow:
1. Select the best-fit runbook and validate symptoms match.
2. Compare with past incident resolutions for proven fixes.
3. Propose low-risk mitigations first (config tuning, scale-up).
4. Require approval before destructive actions (rollback, failover).`,
			serviceName, query, strings.Join(rbLines, "\n"), strings.Join(incLines, "\n"))

		return &mcp.GetPromptResult{
			Description: fmt.Sprintf("Mitigation recommendations for %s", serviceName),
			Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
		}, nil
	})
}
