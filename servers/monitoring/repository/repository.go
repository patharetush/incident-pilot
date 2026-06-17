package repository

import (
	"context"
	"time"
)

// Repository abstracts monitoring data access. Swap implementations (memory,
// Prometheus, Datadog) without changing service or MCP layers.
type Repository interface {
	ListServices(ctx context.Context) ([]Service, error)
	GetService(ctx context.Context, name string) (Service, error)
	GetMetrics(ctx context.Context, service string) ([]MetricPoint, error)
	ListAlerts(ctx context.Context, filter AlertFilter) ([]Alert, error)
	GetAlert(ctx context.Context, id string) (Alert, error)
}

type ServiceStatus string

const (
	StatusHealthy  ServiceStatus = "healthy"
	StatusDegraded ServiceStatus = "degraded"
	StatusCritical ServiceStatus = "critical"
	StatusUnknown  ServiceStatus = "unknown"
)

type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityCritical AlertSeverity = "critical"
)

type Service struct {
	Name        string        `json:"name"`
	Environment string        `json:"environment"`
	Status      ServiceStatus `json:"status"`
	Summary     string        `json:"summary"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type MetricPoint struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

type Alert struct {
	ID          string            `json:"id"`
	Service     string            `json:"service"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Severity    AlertSeverity     `json:"severity"`
	Status      string            `json:"status"`
	StartedAt   time.Time         `json:"started_at"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type AlertFilter struct {
	Severity string
}
