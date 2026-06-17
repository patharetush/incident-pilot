package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/deployments/service"
)

const (
	CatalogURI       = "deployments://catalog/recent"
	DeploymentTemplateURI = "deployments://deployments/{id}"
)

func Register(server *mcp.Server, svc *service.Service) {
	server.AddResource(&mcp.Resource{
		URI:         CatalogURI,
		Name:        "recent-deployments",
		Description: "JSON catalog of recent deployments across services",
		MIMEType:    "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		deployments, err := svc.ListDeployments(ctx, "", 10)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(map[string]any{"deployments": deployments, "count": len(deployments)})
		if err != nil {
			return nil, err
		}
		return jsonResource(CatalogURI, string(data)), nil
	})

	server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: DeploymentTemplateURI,
		Name:        "deployment-detail",
		Description: "Detailed deployment record including changes",
		MIMEType:    "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		id, err := idFromURI(req.Params.URI, "deployments://deployments/")
		if err != nil {
			return nil, err
		}
		dep, err := svc.GetDeployment(ctx, id)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(dep)
		if err != nil {
			return nil, err
		}
		return jsonResource(req.Params.URI, string(data)), nil
	})
}

func idFromURI(uri, prefix string) (string, error) {
	if len(uri) <= len(prefix) || uri[:len(prefix)] != prefix {
		return "", fmt.Errorf("invalid resource URI: %s", uri)
	}
	id := uri[len(prefix):]
	if id == "" {
		return "", fmt.Errorf("id is required in URI")
	}
	return id, nil
}

func jsonResource(uri, text string) *mcp.ReadResourceResult {
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{{
			URI: uri, MIMEType: "application/json", Text: text,
		}},
	}
}
