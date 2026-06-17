package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/deployments/repository"
	"github.com/patharetush/incident-pilot/servers/deployments/service"
)

type ListDeploymentsInput struct {
	Service string `json:"service,omitempty" jsonschema:"Optional service filter (e.g. payment-api)"`
	Limit   int    `json:"limit,omitempty" jsonschema:"Maximum deployments to return (default 20)"`
}

type ListDeploymentsOutput struct {
	Deployments []repository.Deployment `json:"deployments"`
	Count       int                     `json:"count"`
}

func registerListDeployments(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        ListDeployments,
		Description: "List recent deployments, optionally filtered by service",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListDeploymentsInput) (*mcp.CallToolResult, ListDeploymentsOutput, error) {
		deployments, err := svc.ListDeployments(ctx, input.Service, input.Limit)
		if err != nil {
			return toolError(err.Error()), ListDeploymentsOutput{}, nil
		}
		return nil, ListDeploymentsOutput{Deployments: deployments, Count: len(deployments)}, nil
	})
}

type GetDeploymentInput struct {
	DeploymentID string `json:"deployment_id" jsonschema:"Deployment identifier (e.g. dep-001)"`
}

type GetDeploymentOutput struct {
	Deployment repository.Deployment `json:"deployment"`
}

func registerGetDeployment(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        GetDeployment,
		Description: "Get detailed information about a specific deployment",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetDeploymentInput) (*mcp.CallToolResult, GetDeploymentOutput, error) {
		dep, err := svc.GetDeployment(ctx, input.DeploymentID)
		if err != nil {
			return toolError(err.Error()), GetDeploymentOutput{}, nil
		}
		return nil, GetDeploymentOutput{Deployment: dep}, nil
	})
}

type ListRecentChangesInput struct {
	Service      string `json:"service" jsonschema:"Service name (e.g. payment-api)"`
	SinceMinutes int    `json:"since_minutes,omitempty" jsonschema:"Look back window in minutes (default 60)"`
}

type ListRecentChangesOutput struct {
	Service string              `json:"service"`
	Changes []repository.Change `json:"changes"`
	Count   int                 `json:"count"`
}

func registerListRecentChanges(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        ListRecentChanges,
		Description: "List infrastructure and code changes for a service within a recent time window",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListRecentChangesInput) (*mcp.CallToolResult, ListRecentChangesOutput, error) {
		changes, err := svc.ListRecentChanges(ctx, input.Service, input.SinceMinutes)
		if err != nil {
			return toolError(err.Error()), ListRecentChangesOutput{}, nil
		}
		return nil, ListRecentChangesOutput{
			Service: input.Service,
			Changes: changes,
			Count:   len(changes),
		}, nil
	})
}
