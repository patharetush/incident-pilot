package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/logs/service"
)

const (
	PatternsURI      = "logs://patterns/errors"
	ServiceTemplateURI = "logs://service/{name}/recent"
)

func Register(server *mcp.Server, svc *service.Service) {
	server.AddResource(&mcp.Resource{
		URI: PatternsURI, Name: "error-patterns",
		Description: "JSON summary of recurring error patterns", MIMEType: "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		patterns, err := svc.ListErrorPatterns(ctx, "")
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(map[string]any{"patterns": patterns, "count": len(patterns)})
		return jsonResource(PatternsURI, string(data)), nil
	})

	server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: ServiceTemplateURI, Name: "recent-logs",
		Description: "Recent log entries for a service", MIMEType: "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		name, err := nameFromURI(req.Params.URI)
		if err != nil {
			return nil, err
		}
		logs, err := svc.SearchLogs(ctx, name, "", "", 20)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(map[string]any{"service": name, "logs": logs, "count": len(logs)})
		return jsonResource(req.Params.URI, string(data)), nil
	})
}

func nameFromURI(uri string) (string, error) {
	const prefix = "logs://service/"
	const suffix = "/recent"
	if len(uri) <= len(prefix)+len(suffix) {
		return "", fmt.Errorf("invalid log resource URI: %s", uri)
	}
	if uri[:len(prefix)] != prefix || uri[len(uri)-len(suffix):] != suffix {
		return "", fmt.Errorf("invalid log resource URI: %s", uri)
	}
	name := uri[len(prefix) : len(uri)-len(suffix)]
	if name == "" {
		return "", fmt.Errorf("service name required in URI")
	}
	return name, nil
}

func jsonResource(uri, text string) *mcp.ReadResourceResult {
	return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: uri, MIMEType: "application/json", Text: text}}}
}
