package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/knowledge/service"
)

const (
	SearchRunbooks       = "search_runbooks"
	GetRunbook           = "get_runbook"
	SearchPastIncidents  = "search_past_incidents"
)

func Names() []string {
	return []string{SearchRunbooks, GetRunbook, SearchPastIncidents}
}

func Register(server *mcp.Server, svc *service.Service) {
	registerSearchRunbooks(server, svc)
	registerGetRunbook(server, svc)
	registerSearchPastIncidents(server, svc)
}
