package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/logs/service"
)

const AnalyzeLogs = "analyze_logs"

func Register(server *mcp.Server, svc *service.Service) {
	server.AddPrompt(&mcp.Prompt{
		Name: AnalyzeLogs,
		Description: "Structured log analysis workflow for incident evidence gathering",
		Arguments: []*mcp.PromptArgument{{
			Name: "service", Description: "Service to analyze (e.g. payment-api)", Required: true,
		}},
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		serviceName := strings.TrimSpace(req.Params.Arguments["service"])
		if serviceName == "" {
			return nil, fmt.Errorf("service argument is required")
		}

		logs, err := svc.SearchLogs(ctx, serviceName, "error", "", 10)
		if err != nil {
			return nil, err
		}
		patterns, err := svc.ListErrorPatterns(ctx, serviceName)
		if err != nil {
			return nil, err
		}

		var logLines, patternLines []string
		for _, l := range logs {
			logLines = append(logLines, fmt.Sprintf("- [%s] %s", l.Timestamp.Format("15:04:05"), l.Message))
		}
		if len(logLines) == 0 {
			logLines = append(logLines, "- No error logs found")
		}
		for _, p := range patterns {
			patternLines = append(patternLines, fmt.Sprintf("- %s (count=%d): %s", p.Pattern, p.Count, p.Sample))
		}
		if len(patternLines) == 0 {
			patternLines = append(patternLines, "- No error patterns detected")
		}

		text := fmt.Sprintf(`Analyze logs for incident evidence on %q.

Recent error logs:
%s

Error patterns:
%s

Analysis steps:
1. Identify the earliest error signal and trace IDs.
2. Use get_log_context to pull correlated entries.
3. Match patterns against known failure modes.
4. Document timeline evidence before proposing fixes.`,
			serviceName, strings.Join(logLines, "\n"), strings.Join(patternLines, "\n"))

		return &mcp.GetPromptResult{
			Description: fmt.Sprintf("Log analysis for %s", serviceName),
			Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
		}, nil
	})
}
