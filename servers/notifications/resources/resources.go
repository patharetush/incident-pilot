package resources

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/notifications/service"
)

const (
	ChannelsURI = "notifications://channels/catalog"
	PendingURI  = "notifications://pending/approvals"
)

func Register(server *mcp.Server, svc *service.Service) {
	server.AddResource(&mcp.Resource{
		URI: ChannelsURI, Name: "channel-catalog",
		Description: "JSON catalog of notification channels and audiences", MIMEType: "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		channels, err := svc.ListChannels(ctx)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(map[string]any{"channels": channels, "count": len(channels)})
		return jsonResource(ChannelsURI, string(data)), nil
	})

	server.AddResource(&mcp.Resource{
		URI: PendingURI, Name: "pending-approvals",
		Description: "Notifications awaiting human approval", MIMEType: "application/json",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		pending, err := svc.ListPendingApprovals(ctx)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(map[string]any{"notifications": pending, "count": len(pending)})
		return jsonResource(PendingURI, string(data)), nil
	})
}

func jsonResource(uri, text string) *mcp.ReadResourceResult {
	return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: uri, MIMEType: "application/json", Text: text}}}
}
