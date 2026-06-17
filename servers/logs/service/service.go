package service

import (
	"context"
	"fmt"

	"github.com/patharetush/incident-pilot/servers/logs/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SearchLogs(ctx context.Context, service, level, query string, limit int) ([]repository.LogEntry, error) {
	return s.repo.SearchLogs(ctx, repository.LogFilter{
		Service: service,
		Level:   repository.LogLevel(level),
		Query:   query,
		Limit:   limit,
	})
}

func (s *Service) GetLogEntry(ctx context.Context, id string) (repository.LogEntry, error) {
	if id == "" {
		return repository.LogEntry{}, fmt.Errorf("log id is required")
	}
	return s.repo.GetLogEntry(ctx, id)
}

func (s *Service) ListErrorPatterns(ctx context.Context, service string) ([]repository.ErrorPattern, error) {
	return s.repo.ListErrorPatterns(ctx, service)
}

func (s *Service) GetLogContext(ctx context.Context, traceID string, limit int) ([]repository.LogEntry, error) {
	if traceID == "" {
		return nil, fmt.Errorf("trace_id is required")
	}
	if limit <= 0 {
		limit = 20
	}
	all, err := s.repo.SearchLogs(ctx, repository.LogFilter{Limit: 100})
	if err != nil {
		return nil, err
	}
	var out []repository.LogEntry
	for _, entry := range all {
		if entry.TraceID == traceID {
			out = append(out, entry)
		}
	}
	if len(out) == 0 {
		return []repository.LogEntry{}, nil
	}
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
