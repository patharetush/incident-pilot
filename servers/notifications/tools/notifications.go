package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/notifications/repository"
	"github.com/patharetush/incident-pilot/servers/notifications/service"
)

type ListChannelsOutput struct {
	Channels []repository.Channel `json:"channels"`
	Count    int                  `json:"count"`
}

func registerListChannels(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: ListChannels, Description: "List available notification channels (Slack, PagerDuty, email)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListChannelsOutput, error) {
		channels, err := svc.ListChannels(ctx)
		if err != nil {
			return toolError(err.Error()), ListChannelsOutput{}, nil
		}
		return nil, ListChannelsOutput{Channels: channels, Count: len(channels)}, nil
	})
}

type ListNotificationsInput struct {
	Status  string `json:"status,omitempty" jsonschema:"Filter by status: sent, queued, pending_approval, failed"`
	Service string `json:"service,omitempty" jsonschema:"Optional service filter"`
	Limit   int    `json:"limit,omitempty" jsonschema:"Max results (default 50)"`
}

type ListNotificationsOutput struct {
	Notifications []repository.Notification `json:"notifications"`
	Count         int                       `json:"count"`
}

func registerListNotifications(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: ListNotifications, Description: "List sent and queued notifications for an incident",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListNotificationsInput) (*mcp.CallToolResult, ListNotificationsOutput, error) {
		notifications, err := svc.ListNotifications(ctx, input.Status, input.Service, input.Limit)
		if err != nil {
			return toolError(err.Error()), ListNotificationsOutput{}, nil
		}
		return nil, ListNotificationsOutput{Notifications: notifications, Count: len(notifications)}, nil
	})
}

type GetNotificationInput struct {
	NotificationID string `json:"notification_id" jsonschema:"Notification identifier (e.g. notif-001)"`
}

type GetNotificationOutput struct {
	Notification repository.Notification `json:"notification"`
}

func registerGetNotification(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: GetNotification, Description: "Get full details of a notification by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetNotificationInput) (*mcp.CallToolResult, GetNotificationOutput, error) {
		n, err := svc.GetNotification(ctx, input.NotificationID)
		if err != nil {
			return toolError(err.Error()), GetNotificationOutput{}, nil
		}
		return nil, GetNotificationOutput{Notification: n}, nil
	})
}

type QueueNotificationInput struct {
	ChannelID        string `json:"channel_id" jsonschema:"Target channel ID (e.g. ch-slack-incidents)"`
	Subject          string `json:"subject" jsonschema:"Notification subject line"`
	Body             string `json:"body" jsonschema:"Notification body content"`
	IncidentID       string `json:"incident_id,omitempty" jsonschema:"Related incident ID"`
	Service          string `json:"service,omitempty" jsonschema:"Affected service"`
	RequestedBy      string `json:"requested_by,omitempty" jsonschema:"Actor requesting the notification"`
	RequiresApproval bool   `json:"requires_approval,omitempty" jsonschema:"If true, holds for human approval before sending"`
}

type QueueNotificationOutput struct {
	Notification repository.Notification `json:"notification"`
}

func registerQueueNotification(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: QueueNotification,
		Description: "Queue a notification for delivery; high-impact channels require approval",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input QueueNotificationInput) (*mcp.CallToolResult, QueueNotificationOutput, error) {
		n, err := svc.QueueNotification(ctx, input.ChannelID, input.Subject, input.Body,
			input.IncidentID, input.Service, input.RequestedBy, input.RequiresApproval)
		if err != nil {
			return toolError(err.Error()), QueueNotificationOutput{}, nil
		}
		return nil, QueueNotificationOutput{Notification: n}, nil
	})
}

type ListPendingApprovalsOutput struct {
	Notifications []repository.Notification `json:"notifications"`
	Count         int                     `json:"count"`
}

func registerListPendingApprovals(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: ListPendingApprovals, Description: "List notifications awaiting human approval before send",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListPendingApprovalsOutput, error) {
		notifications, err := svc.ListPendingApprovals(ctx)
		if err != nil {
			return toolError(err.Error()), ListPendingApprovalsOutput{}, nil
		}
		return nil, ListPendingApprovalsOutput{Notifications: notifications, Count: len(notifications)}, nil
	})
}

type ApproveNotificationInput struct {
	NotificationID string `json:"notification_id" jsonschema:"Pending notification to approve"`
	Approver       string `json:"approver" jsonschema:"Name or ID of the approving operator"`
}

type ApproveNotificationOutput struct {
	Notification repository.Notification `json:"notification"`
}

func registerApproveNotification(server *mcp.Server, svc *service.Service) {
	mcp.AddTool(server, &mcp.Tool{
		Name: ApproveNotification, Description: "Approve and send a pending notification (audit trail recorded)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ApproveNotificationInput) (*mcp.CallToolResult, ApproveNotificationOutput, error) {
		n, err := svc.ApproveNotification(ctx, input.NotificationID, input.Approver)
		if err != nil {
			return toolError(err.Error()), ApproveNotificationOutput{}, nil
		}
		return nil, ApproveNotificationOutput{Notification: n}, nil
	})
}
