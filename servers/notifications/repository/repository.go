package repository

import (
	"context"
	"time"
)

type Repository interface {
	ListChannels(ctx context.Context) ([]Channel, error)
	ListNotifications(ctx context.Context, filter NotificationFilter) ([]Notification, error)
	GetNotification(ctx context.Context, id string) (Notification, error)
	QueueNotification(ctx context.Context, req QueueRequest) (Notification, error)
	ListPendingApprovals(ctx context.Context) ([]Notification, error)
	ApproveNotification(ctx context.Context, id, approver string) (Notification, error)
}

type ChannelType string

const (
	ChannelSlack     ChannelType = "slack"
	ChannelPagerDuty ChannelType = "pagerduty"
	ChannelEmail     ChannelType = "email"
)

type NotificationStatus string

const (
	StatusSent    NotificationStatus = "sent"
	StatusQueued  NotificationStatus = "queued"
	StatusPending NotificationStatus = "pending_approval"
	StatusFailed  NotificationStatus = "failed"
)

type Channel struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        ChannelType `json:"type"`
	Description string      `json:"description"`
	Audience    string      `json:"audience"`
}

type Notification struct {
	ID          string             `json:"id"`
	ChannelID   string             `json:"channel_id"`
	Subject     string             `json:"subject"`
	Body        string             `json:"body"`
	Status      NotificationStatus `json:"status"`
	IncidentID  string             `json:"incident_id,omitempty"`
	Service     string             `json:"service,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	SentAt      *time.Time         `json:"sent_at,omitempty"`
	RequestedBy string             `json:"requested_by,omitempty"`
	ApprovedBy  string             `json:"approved_by,omitempty"`
}

type QueueRequest struct {
	ChannelID   string
	Subject     string
	Body        string
	IncidentID  string
	Service     string
	RequestedBy string
	RequiresApproval bool
}

type NotificationFilter struct {
	Status  NotificationStatus
	Service string
	Limit   int
}
