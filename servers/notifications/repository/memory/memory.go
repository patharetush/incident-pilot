package memory

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/patharetush/incident-pilot/servers/notifications/repository"
)

type Repository struct {
	mu            sync.Mutex
	channels      map[string]repository.Channel
	notifications map[string]repository.Notification
	nextID        int
}

func New() *Repository {
	now := time.Now().UTC().Truncate(time.Second)
	sentAt := now.Add(-10 * time.Minute)
	channels := map[string]repository.Channel{
		"ch-slack-incidents": {
			ID: "ch-slack-incidents", Name: "#incidents",
			Type: repository.ChannelSlack, Description: "Primary incident coordination channel",
			Audience: "on-call, incident commanders",
		},
		"ch-slack-payments": {
			ID: "ch-slack-payments", Name: "#team-payments",
			Type: repository.ChannelSlack, Description: "Payments team updates",
			Audience: "payments engineering",
		},
		"ch-pagerduty": {
			ID: "ch-pagerduty", Name: "PagerDuty - Production",
			Type: repository.ChannelPagerDuty, Description: "Production on-call paging",
			Audience: "primary on-call",
		},
		"ch-email-leads": {
			ID: "ch-email-leads", Name: "Engineering Leads",
			Type: repository.ChannelEmail, Description: "Executive incident summaries",
			Audience: "engineering leadership",
		},
	}
	notifications := map[string]repository.Notification{
		"notif-001": {
			ID: "notif-001", ChannelID: "ch-pagerduty",
			Subject: "SEV-2: payment-api elevated error rate",
			Body:    "5xx error rate on payment-api exceeded 5% SLO. Investigating correlation with deployment v2.14.0.",
			Status: repository.StatusSent, IncidentID: "inc-current",
			Service: "payment-api", CreatedAt: now.Add(-11 * time.Minute),
			SentAt: &sentAt, RequestedBy: "incident-bot", ApprovedBy: "auto",
		},
		"notif-002": {
			ID: "notif-002", ChannelID: "ch-slack-incidents",
			Subject: "Incident update: payment-api degraded",
			Body:    "Status: investigating. Alerts: HighErrorRate, LatencySLOBreach. Next update in 15 minutes.",
			Status: repository.StatusSent, IncidentID: "inc-current",
			Service: "payment-api", CreatedAt: now.Add(-8 * time.Minute),
			SentAt: ptrTime(now.Add(-8 * time.Minute)), RequestedBy: "incident-bot", ApprovedBy: "oncall-alice",
		},
		"notif-003": {
			ID: "notif-003", ChannelID: "ch-email-leads",
			Subject: "Draft: Executive summary - payment-api incident",
			Body:    "Customer impact: elevated payment failures. Root cause under investigation. ETA for mitigation TBD.",
			Status: repository.StatusPending, IncidentID: "inc-current",
			Service: "payment-api", CreatedAt: now.Add(-3 * time.Minute),
			RequestedBy: "incident-bot",
		},
	}
	return &Repository{
		channels: channels, notifications: notifications, nextID: 100,
	}
}

func ptrTime(t time.Time) *time.Time { return &t }

func (r *Repository) ListChannels(ctx context.Context) ([]repository.Channel, error) {
	out := make([]repository.Channel, 0, len(r.channels))
	for _, ch := range r.channels {
		out = append(out, ch)
	}
	return out, nil
}

func (r *Repository) ListNotifications(ctx context.Context, filter repository.NotificationFilter) ([]repository.Notification, error) {
	out := make([]repository.Notification, 0)
	for _, n := range r.notifications {
		if filter.Status != "" && n.Status != filter.Status {
			continue
		}
		if filter.Service != "" && !strings.EqualFold(n.Service, filter.Service) {
			continue
		}
		out = append(out, n)
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *Repository) GetNotification(ctx context.Context, id string) (repository.Notification, error) {
	n, ok := r.notifications[strings.ToLower(id)]
	if !ok {
		return repository.Notification{}, fmt.Errorf("notification %q not found", id)
	}
	return n, nil
}

func (r *Repository) QueueNotification(ctx context.Context, req repository.QueueRequest) (repository.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.channels[req.ChannelID]; !ok {
		return repository.Notification{}, fmt.Errorf("channel %q not found", req.ChannelID)
	}
	if req.Subject == "" || req.Body == "" {
		return repository.Notification{}, fmt.Errorf("subject and body are required")
	}

	r.nextID++
	id := fmt.Sprintf("notif-%03d", r.nextID)
	status := repository.StatusQueued
	if req.RequiresApproval {
		status = repository.StatusPending
	}

	n := repository.Notification{
		ID: id, ChannelID: req.ChannelID, Subject: req.Subject, Body: req.Body,
		Status: status, IncidentID: req.IncidentID, Service: req.Service,
		CreatedAt: time.Now().UTC().Truncate(time.Second), RequestedBy: req.RequestedBy,
	}
	r.notifications[id] = n
	return n, nil
}

func (r *Repository) ListPendingApprovals(ctx context.Context) ([]repository.Notification, error) {
	return r.ListNotifications(ctx, repository.NotificationFilter{Status: repository.StatusPending, Limit: 50})
}

func (r *Repository) ApproveNotification(ctx context.Context, id, approver string) (repository.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := strings.ToLower(id)
	n, ok := r.notifications[key]
	if !ok {
		return repository.Notification{}, fmt.Errorf("notification %q not found", id)
	}
	if n.Status != repository.StatusPending {
		return repository.Notification{}, fmt.Errorf("notification %q is not pending approval (status: %s)", id, n.Status)
	}
	if approver == "" {
		return repository.Notification{}, fmt.Errorf("approver is required")
	}

	now := time.Now().UTC().Truncate(time.Second)
	n.Status = repository.StatusSent
	n.ApprovedBy = approver
	n.SentAt = &now
	r.notifications[key] = n
	return n, nil
}
