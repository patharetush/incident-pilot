package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/knowledge/service"
)

const (
	CatalogURI         = "knowledge://runbooks/catalog"
	RunbookTemplateURI = "knowledge://runbooks/{id}"
)

func Register(server *mcp.Server, svc *service.Service) {
	server.AddResource(&mcp.Resource{
		URI: CatalogURI, Name: "runbook-catalog",
		Description: "JSON catalog of available runbooks", MIMEType: "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		runbooks, err := svc.SearchRunbooks(ctx, "", "", 50)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(map[string]any{"runbooks": runbooks, "count": len(runbooks)})
		return jsonResource(CatalogURI, string(data)), nil
	})

	server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: RunbookTemplateURI, Name: "runbook-detail",
		Description: "Full runbook with remediation steps", MIMEType: "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		id, err := idFromURI(req.Params.URI, "knowledge://runbooks/")
		if err != nil {
			return nil, err
		}
		rb, err := svc.GetRunbook(ctx, id)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(rb)
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
	return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: uri, MIMEType: "application/json", Text: text}}}
}
