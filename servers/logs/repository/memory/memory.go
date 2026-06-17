package memory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/patharetush/incident-pilot/servers/logs/repository"
)

type Repository struct {
	logs     []repository.LogEntry
	patterns []repository.ErrorPattern
}

func New() *Repository {
	now := time.Now().UTC().Truncate(time.Second)
	logs := []repository.LogEntry{
		{
			ID: "log-001", Service: "payment-api", Level: repository.LevelError,
			Message: "upstream timeout after 800ms calling payment-gateway",
			Timestamp: now.Add(-14 * time.Minute), TraceID: "trace-abc123",
			Metadata: map[string]string{"endpoint": "/v1/charge", "status": "504"},
		},
		{
			ID: "log-002", Service: "payment-api", Level: repository.LevelError,
			Message: "retry exhausted for payment-gateway request",
			Timestamp: now.Add(-13 * time.Minute), TraceID: "trace-abc123",
		},
		{
			ID: "log-003", Service: "payment-api", Level: repository.LevelWarn,
			Message: "elevated 5xx rate detected on /v1/charge",
			Timestamp: now.Add(-12 * time.Minute),
		},
		{
			ID: "log-004", Service: "inventory-db", Level: repository.LevelError,
			Message: "connection pool exhausted: timeout waiting for connection",
			Timestamp: now.Add(-9 * time.Minute), TraceID: "trace-db789",
		},
		{
			ID: "log-005", Service: "user-service", Level: repository.LevelInfo,
			Message: "health check passed",
			Timestamp: now.Add(-2 * time.Minute),
		},
	}
	patterns := []repository.ErrorPattern{
		{Pattern: "upstream timeout", Count: 342, Service: "payment-api", Sample: logs[0].Message},
		{Pattern: "retry exhausted", Count: 128, Service: "payment-api", Sample: logs[1].Message},
		{Pattern: "connection pool exhausted", Count: 89, Service: "inventory-db", Sample: logs[3].Message},
	}
	return &Repository{logs: logs, patterns: patterns}
}

func (r *Repository) SearchLogs(ctx context.Context, filter repository.LogFilter) ([]repository.LogEntry, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	var out []repository.LogEntry
	for _, entry := range r.logs {
		if filter.Service != "" && !strings.EqualFold(entry.Service, filter.Service) {
			continue
		}
		if filter.Level != "" && entry.Level != filter.Level {
			continue
		}
		if filter.Query != "" && !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(filter.Query)) {
			continue
		}
		out = append(out, entry)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (r *Repository) GetLogEntry(ctx context.Context, id string) (repository.LogEntry, error) {
	for _, entry := range r.logs {
		if strings.EqualFold(entry.ID, id) {
			return entry, nil
		}
	}
	return repository.LogEntry{}, fmt.Errorf("log entry %q not found", id)
}

func (r *Repository) ListErrorPatterns(ctx context.Context, service string) ([]repository.ErrorPattern, error) {
	var out []repository.ErrorPattern
	for _, p := range r.patterns {
		if service != "" && !strings.EqualFold(p.Service, service) {
			continue
		}
		out = append(out, p)
	}
	return out, nil
}
