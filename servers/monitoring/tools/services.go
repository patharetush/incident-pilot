package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring/repository"
	"github.com/patharetush/incident-pilot/servers/monitoring/service"
)

type ListServicesOutput struct {
	Services []repository.Service `json:"services" jsonschema:"Monitored services and their current health status"`
	Count    int                  `json:"count" jsonschema:"Number of services returned"`
}

func registerListServices(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        ListServices,
		Description: "List all monitored services with current health status and summary",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListServicesOutput, error) {
		services, err := svc.ListServices(ctx)
		if err != nil {
			return toolError(err.Error()), ListServicesOutput{}, nil
		}
		return nil, ListServicesOutput{Services: services, Count: len(services)}, nil
	})
}
