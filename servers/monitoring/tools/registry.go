package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring/service"
)

const (
	ListServices      = "list_services"
	GetServiceMetrics = "get_service_metrics"
	ListAlerts        = "list_alerts"
	GetAlert          = "get_alert"
)

// Names returns all registered tool names.
func Names() []string {
	return []string{ListServices, GetServiceMetrics, ListAlerts, GetAlert}
}

// Register attaches all monitoring tools to the MCP server.
func Register(server *mcp.Server, svc *service.Service) {
	registerListServices(server, svc)
	registerGetServiceMetrics(server, svc)
	registerListAlerts(server, svc)
	registerGetAlert(server, svc)
}
