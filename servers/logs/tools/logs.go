package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/logs/repository"
	"github.com/patharetush/incident-pilot/servers/logs/service"
)

type SearchLogsInput struct {
	Service string `json:"service,omitempty" jsonschema:"Filter by service name"`
	Level   string `json:"level,omitempty" jsonschema:"Filter by level: debug, info, warn, error"`
	Query   string `json:"query,omitempty" jsonschema:"Substring match on log message"`
	Limit   int    `json:"limit,omitempty" jsonschema:"Max entries to return (default 50)"`
}

type SearchLogsOutput struct {
	Logs  []repository.LogEntry `json:"logs"`
	Count int                   `json:"count"`
}

func registerSearchLogs(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: SearchLogs, Description: "Search log entries by service, level, and message query",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SearchLogsInput) (*mcp.CallToolResult, SearchLogsOutput, error) {
		logs, err := svc.SearchLogs(ctx, input.Service, input.Level, input.Query, input.Limit)
		if err != nil {
			return toolError(err.Error()), SearchLogsOutput{}, nil
		}
		return nil, SearchLogsOutput{Logs: logs, Count: len(logs)}, nil
	})
}

type GetLogEntryInput struct {
	LogID string `json:"log_id" jsonschema:"Log entry identifier (e.g. log-001)"`
}

type GetLogEntryOutput struct {
	Log repository.LogEntry `json:"log"`
}

func registerGetLogEntry(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: GetLogEntry, Description: "Get a specific log entry by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetLogEntryInput) (*mcp.CallToolResult, GetLogEntryOutput, error) {
		log, err := svc.GetLogEntry(ctx, input.LogID)
		if err != nil {
			return toolError(err.Error()), GetLogEntryOutput{}, nil
		}
		return nil, GetLogEntryOutput{Log: log}, nil
	})
}

type GetLogContextInput struct {
	TraceID string `json:"trace_id" jsonschema:"Distributed trace ID to fetch related log entries"`
	Limit   int    `json:"limit,omitempty" jsonschema:"Max related entries (default 20)"`
}

type GetLogContextOutput struct {
	TraceID string                `json:"trace_id"`
	Logs    []repository.LogEntry `json:"logs"`
	Count   int                   `json:"count"`
}

func registerGetLogContext(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: GetLogContext, Description: "Get log entries sharing the same trace ID for correlation",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetLogContextInput) (*mcp.CallToolResult, GetLogContextOutput, error) {
		logs, err := svc.GetLogContext(ctx, input.TraceID, input.Limit)
		if err != nil {
			return toolError(err.Error()), GetLogContextOutput{}, nil
		}
		return nil, GetLogContextOutput{TraceID: input.TraceID, Logs: logs, Count: len(logs)}, nil
	})
}

type ListErrorPatternsInput struct {
	Service string `json:"service,omitempty" jsonschema:"Optional service filter"`
}

type ListErrorPatternsOutput struct {
	Patterns []repository.ErrorPattern `json:"patterns"`
	Count    int                       `json:"count"`
}

func registerListErrorPatterns(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: ListErrorPatterns, Description: "List recurring error patterns detected in recent logs",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListErrorPatternsInput) (*mcp.CallToolResult, ListErrorPatternsOutput, error) {
		patterns, err := svc.ListErrorPatterns(ctx, input.Service)
		if err != nil {
			return toolError(err.Error()), ListErrorPatternsOutput{}, nil
		}
		return nil, ListErrorPatternsOutput{Patterns: patterns, Count: len(patterns)}, nil
	})
}
