package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/logs/service"
)

const (
	SearchLogs         = "search_logs"
	GetLogEntry        = "get_log_entry"
	GetLogContext      = "get_log_context"
	ListErrorPatterns  = "list_error_patterns"
)

func Names() []string {
	return []string{SearchLogs, GetLogEntry, GetLogContext, ListErrorPatterns}
}

func Register(server *mcp.Server, svc *service.Service) {
	registerSearchLogs(server, svc)
	registerGetLogEntry(server, svc)
	registerGetLogContext(server, svc)
	registerListErrorPatterns(server, svc)
}
