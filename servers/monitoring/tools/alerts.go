package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring/repository"
	"github.com/patharetush/incident-pilot/servers/monitoring/service"
)

type ListAlertsInput struct {
	Severity string `json:"severity,omitempty" jsonschema:"Optional filter: info, warning, or critical"`
}

type ListAlertsOutput struct {
	Alerts []repository.Alert `json:"alerts"`
	Count  int                `json:"count"`
}

type GetAlertInput struct {
	AlertID string `json:"alert_id" jsonschema:"Alert identifier (e.g. alert-001)"`
}

type GetAlertOutput struct {
	Alert repository.Alert `json:"alert"`
}

func registerListAlerts(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        ListAlerts,
		Description: "List active monitoring alerts, optionally filtered by severity",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListAlertsInput) (*mcp.CallToolResult, ListAlertsOutput, error) {
		alerts, err := svc.ListAlerts(ctx, input.Severity)
		if err != nil {
			return toolError(err.Error()), ListAlertsOutput{}, nil
		}
		return nil, ListAlertsOutput{Alerts: alerts, Count: len(alerts)}, nil
	})
}

func registerGetAlert(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        GetAlert,
		Description: "Get detailed information about a specific alert by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetAlertInput) (*mcp.CallToolResult, GetAlertOutput, error) {
		alert, err := svc.GetAlert(ctx, input.AlertID)
		if err != nil {
			return toolError(err.Error()), GetAlertOutput{}, nil
		}
		return nil, GetAlertOutput{Alert: alert}, nil
	})
}
