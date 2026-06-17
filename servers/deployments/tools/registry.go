package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/deployments/service"
)

const (
	ListDeployments    = "list_deployments"
	GetDeployment      = "get_deployment"
	ListRecentChanges  = "list_recent_changes"
)

func Names() []string {
	return []string{ListDeployments, GetDeployment, ListRecentChanges}
}

func Register(server *mcp.Server, svc *service.Service) {
	registerListDeployments(server, svc)
	registerGetDeployment(server, svc)
	registerListRecentChanges(server, svc)
}
