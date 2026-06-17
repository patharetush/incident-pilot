package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/monitoring/service"
)

const (
	CatalogURI        = "monitoring://catalog/services"
	ServiceTemplateURI = "monitoring://services/{name}"
)

// Register attaches monitoring MCP resources to the server.
func Register(server *mcp.Server, svc *service.Service) {
	server.AddResource(&mcp.Resource{
		URI:         CatalogURI,
		Name:        "service-catalog",
		Description: "JSON catalog of all monitored services and their health status",
		MIMEType:    "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		services, err := svc.ListServices(ctx)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(map[string]any{
			"services": services,
			"count":    len(services),
		})
		if err != nil {
			return nil, err
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      CatalogURI,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	})

	server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: ServiceTemplateURI,
		Name:        "service-detail",
		Description: "Detailed status for a specific monitored service",
		MIMEType:    "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		name, err := serviceNameFromURI(req.Params.URI)
		if err != nil {
			return nil, err
		}
		svcDetail, err := svc.GetService(ctx, name)
		if err != nil {
			return nil, err
		}
		_, metricPoints, err := svc.GetServiceMetrics(ctx, name)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(map[string]any{
			"service": svcDetail,
			"metrics": metricPoints,
		})
		if err != nil {
			return nil, err
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	})
}

func serviceNameFromURI(uri string) (string, error) {
	const prefix = "monitoring://services/"
	if len(uri) <= len(prefix) || uri[:len(prefix)] != prefix {
		return "", fmt.Errorf("invalid service resource URI: %s", uri)
	}
	name := uri[len(prefix):]
	if name == "" {
		return "", fmt.Errorf("service name is required in URI")
	}
	return name, nil
}
