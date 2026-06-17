package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/servers/notifications/service"
)

const (
	ListChannels         = "list_channels"
	ListNotifications    = "list_notifications"
	GetNotification      = "get_notification"
	QueueNotification    = "queue_notification"
	ListPendingApprovals = "list_pending_approvals"
	ApproveNotification  = "approve_notification"
)

func Names() []string {
	return []string{
		ListChannels, ListNotifications, GetNotification,
		QueueNotification, ListPendingApprovals, ApproveNotification,
	}
}

func Register(server *mcp.Server, svc *service.Service) {
	registerListChannels(server, svc)
	registerListNotifications(server, svc)
	registerGetNotification(server, svc)
	registerQueueNotification(server, svc)
	registerListPendingApprovals(server, svc)
	registerApproveNotification(server, svc)
}
