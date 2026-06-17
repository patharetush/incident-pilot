package memory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/patharetush/incident-pilot/servers/monitoring/repository"
)

// Repository is an in-memory monitoring backend with seeded incident data.
type Repository struct {
	services map[string]repository.Service
	metrics  map[string][]repository.MetricPoint
	alerts   map[string]repository.Alert
}

func New() *Repository {
	now := time.Now().UTC().Truncate(time.Second)
	services := map[string]repository.Service{
		"payment-api": {
			Name:        "payment-api",
			Environment: "production",
			Status:      repository.StatusDegraded,
			Summary:     "Elevated 5xx error rate after recent deployment",
			UpdatedAt:   now,
		},
		"user-service": {
			Name:        "user-service",
			Environment: "production",
			Status:      repository.StatusHealthy,
			Summary:     "All SLOs within target",
			UpdatedAt:   now,
		},
		"inventory-db": {
			Name:        "inventory-db",
			Environment: "production",
			Status:      repository.StatusCritical,
			Summary:     "Connection pool exhausted; queries timing out",
			UpdatedAt:   now,
		},
	}

	metrics := map[string][]repository.MetricPoint{
		"payment-api": {
			{Name: "error_rate", Value: 8.4, Unit: "percent", Timestamp: now},
			{Name: "request_rate", Value: 1240, Unit: "req/s", Timestamp: now},
			{Name: "p99_latency_ms", Value: 890, Unit: "ms", Timestamp: now},
			{Name: "cpu_utilization", Value: 72, Unit: "percent", Timestamp: now},
		},
		"user-service": {
			{Name: "error_rate", Value: 0.2, Unit: "percent", Timestamp: now},
			{Name: "request_rate", Value: 3400, Unit: "req/s", Timestamp: now},
			{Name: "p99_latency_ms", Value: 120, Unit: "ms", Timestamp: now},
			{Name: "cpu_utilization", Value: 41, Unit: "percent", Timestamp: now},
		},
		"inventory-db": {
			{Name: "active_connections", Value: 500, Unit: "connections", Timestamp: now},
			{Name: "max_connections", Value: 500, Unit: "connections", Timestamp: now},
			{Name: "query_timeout_rate", Value: 34.7, Unit: "percent", Timestamp: now},
			{Name: "replication_lag_ms", Value: 2100, Unit: "ms", Timestamp: now},
		},
	}

	alerts := map[string]repository.Alert{
		"alert-001": {
			ID:          "alert-001",
			Service:     "payment-api",
			Title:       "HighErrorRate",
			Description: "5xx error rate exceeded 5% for 10 minutes (current: 8.4%)",
			Severity:    repository.SeverityCritical,
			Status:      "firing",
			StartedAt:   now.Add(-12 * time.Minute),
			Labels: map[string]string{
				"team":       "payments",
				"deployment": "payment-api-v2.14.0",
			},
		},
		"alert-002": {
			ID:          "alert-002",
			Service:     "inventory-db",
			Title:       "ConnectionPoolExhausted",
			Description: "Database connection pool at 100% capacity with rising query timeouts",
			Severity:    repository.SeverityCritical,
			Status:      "firing",
			StartedAt:   now.Add(-8 * time.Minute),
			Labels: map[string]string{
				"team": "platform",
			},
		},
		"alert-003": {
			ID:          "alert-003",
			Service:     "payment-api",
			Title:       "LatencySLOBreach",
			Description: "P99 latency above 500ms SLO threshold",
			Severity:    repository.SeverityWarning,
			Status:      "firing",
			StartedAt:   now.Add(-15 * time.Minute),
			Labels: map[string]string{
				"team": "payments",
			},
		},
	}

	return &Repository{
		services: services,
		metrics:  metrics,
		alerts:   alerts,
	}
}

func (r *Repository) ListServices(ctx context.Context) ([]repository.Service, error) {
	out := make([]repository.Service, 0, len(r.services))
	for _, svc := range r.services {
		out = append(out, svc)
	}
	return out, nil
}

func (r *Repository) GetService(ctx context.Context, name string) (repository.Service, error) {
	svc, ok := r.services[strings.ToLower(name)]
	if !ok {
		return repository.Service{}, fmt.Errorf("service %q not found", name)
	}
	return svc, nil
}

func (r *Repository) GetMetrics(ctx context.Context, service string) ([]repository.MetricPoint, error) {
	key := strings.ToLower(service)
	if _, ok := r.services[key]; !ok {
		return nil, fmt.Errorf("service %q not found", service)
	}
	points, ok := r.metrics[key]
	if !ok {
		return []repository.MetricPoint{}, nil
	}
	out := make([]repository.MetricPoint, len(points))
	copy(out, points)
	return out, nil
}

func (r *Repository) ListAlerts(ctx context.Context, filter repository.AlertFilter) ([]repository.Alert, error) {
	out := make([]repository.Alert, 0, len(r.alerts))
	for _, alert := range r.alerts {
		if filter.Severity != "" && !strings.EqualFold(string(alert.Severity), filter.Severity) {
			continue
		}
		out = append(out, alert)
	}
	return out, nil
}

func (r *Repository) GetAlert(ctx context.Context, id string) (repository.Alert, error) {
	alert, ok := r.alerts[strings.ToLower(id)]
	if !ok {
		return repository.Alert{}, fmt.Errorf("alert %q not found", id)
	}
	return alert, nil
}
