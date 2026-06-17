package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring/repository"
	"github.com/patharetush/incident-pilot/servers/monitoring/service"
)

type GetServiceMetricsInput struct {
	Service string `json:"service" jsonschema:"Service name (e.g. payment-api, user-service, inventory-db)"`
}

type GetServiceMetricsOutput struct {
	Service string                  `json:"service"`
	Metrics []repository.MetricPoint `json:"metrics"`
}

func registerGetServiceMetrics(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        GetServiceMetrics,
		Description: "Get current metrics for a monitored service (error rate, latency, resource utilization)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetServiceMetricsInput) (*mcp.CallToolResult, GetServiceMetricsOutput, error) {
		serviceName, metrics, err := svc.GetServiceMetrics(ctx, input.Service)
		if err != nil {
			return toolError(err.Error()), GetServiceMetricsOutput{}, nil
		}
		return nil, GetServiceMetricsOutput{
			Service: serviceName,
			Metrics: metrics,
		}, nil
	})
}
