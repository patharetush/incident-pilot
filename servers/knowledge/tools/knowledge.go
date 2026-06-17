package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/knowledge/repository"
	"github.com/patharetush/incident-pilot/servers/knowledge/service"
)

type SearchRunbooksInput struct {
	Query   string `json:"query,omitempty" jsonschema:"Search query matched against title, summary, and tags"`
	Service string `json:"service,omitempty" jsonschema:"Optional service filter"`
	Limit   int    `json:"limit,omitempty" jsonschema:"Max results (default 20)"`
}

type SearchRunbooksOutput struct {
	Runbooks []repository.Runbook `json:"runbooks"`
	Count    int                  `json:"count"`
}

func registerSearchRunbooks(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: SearchRunbooks, Description: "Search operational runbooks by keyword and service",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SearchRunbooksInput) (*mcp.CallToolResult, SearchRunbooksOutput, error) {
		runbooks, err := svc.SearchRunbooks(ctx, input.Query, input.Service, input.Limit)
		if err != nil {
			return toolError(err.Error()), SearchRunbooksOutput{}, nil
		}
		return nil, SearchRunbooksOutput{Runbooks: runbooks, Count: len(runbooks)}, nil
	})
}

type GetRunbookInput struct {
	RunbookID string `json:"runbook_id" jsonschema:"Runbook identifier (e.g. rb-001)"`
}

type GetRunbookOutput struct {
	Runbook repository.Runbook `json:"runbook"`
}

func registerGetRunbook(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: GetRunbook, Description: "Get a runbook with full remediation steps",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetRunbookInput) (*mcp.CallToolResult, GetRunbookOutput, error) {
		rb, err := svc.GetRunbook(ctx, input.RunbookID)
		if err != nil {
			return toolError(err.Error()), GetRunbookOutput{}, nil
		}
		return nil, GetRunbookOutput{Runbook: rb}, nil
	})
}

type SearchPastIncidentsInput struct {
	Query   string `json:"query,omitempty" jsonschema:"Search query for title, summary, resolution, tags"`
	Service string `json:"service,omitempty" jsonschema:"Optional service filter"`
	Limit   int    `json:"limit,omitempty" jsonschema:"Max results (default 20)"`
}

type SearchPastIncidentsOutput struct {
	Incidents []repository.PastIncident `json:"incidents"`
	Count     int                       `json:"count"`
}

func registerSearchPastIncidents(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: SearchPastIncidents, Description: "Search historical incidents for similar past events and resolutions",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SearchPastIncidentsInput) (*mcp.CallToolResult, SearchPastIncidentsOutput, error) {
		incidents, err := svc.SearchPastIncidents(ctx, input.Query, input.Service, input.Limit)
		if err != nil {
			return toolError(err.Error()), SearchPastIncidentsOutput{}, nil
		}
		return nil, SearchPastIncidentsOutput{Incidents: incidents, Count: len(incidents)}, nil
	})
}
