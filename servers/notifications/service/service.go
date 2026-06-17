package service

import (
	"context"
	"fmt"

	"github.com/patharetush/incident-pilot/servers/notifications/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListChannels(ctx context.Context) ([]repository.Channel, error) {
	return s.repo.ListChannels(ctx)
}

func (s *Service) ListNotifications(ctx context.Context, status, service string, limit int) ([]repository.Notification, error) {
	return s.repo.ListNotifications(ctx, repository.NotificationFilter{
		Status: repository.NotificationStatus(status),
		Service: service,
		Limit:  limit,
	})
}

func (s *Service) GetNotification(ctx context.Context, id string) (repository.Notification, error) {
	if id == "" {
		return repository.Notification{}, fmt.Errorf("notification id is required")
	}
	return s.repo.GetNotification(ctx, id)
}

func (s *Service) QueueNotification(ctx context.Context, channelID, subject, body, incidentID, service, requestedBy string, requiresApproval bool) (repository.Notification, error) {
	if channelID == "" {
		return repository.Notification{}, fmt.Errorf("channel_id is required")
	}
	if requestedBy == "" {
		requestedBy = "orchestrator"
	}
	return s.repo.QueueNotification(ctx, repository.QueueRequest{
		ChannelID: channelID, Subject: subject, Body: body,
		IncidentID: incidentID, Service: service, RequestedBy: requestedBy,
		RequiresApproval: requiresApproval,
	})
}

func (s *Service) ListPendingApprovals(ctx context.Context) ([]repository.Notification, error) {
	return s.repo.ListPendingApprovals(ctx)
}

func (s *Service) ApproveNotification(ctx context.Context, id, approver string) (repository.Notification, error) {
	if id == "" {
		return repository.Notification{}, fmt.Errorf("notification id is required")
	}
	if approver == "" {
		return repository.Notification{}, fmt.Errorf("approver is required")
	}
	return s.repo.ApproveNotification(ctx, id, approver)
}
