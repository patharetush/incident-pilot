package service

import (
	"context"
	"fmt"

	"github.com/patharetush/incident-pilot/servers/knowledge/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SearchRunbooks(ctx context.Context, query, service string, limit int) ([]repository.Runbook, error) {
	return s.repo.SearchRunbooks(ctx, repository.RunbookFilter{
		Query: query, Service: service, Limit: limit,
	})
}

func (s *Service) GetRunbook(ctx context.Context, id string) (repository.Runbook, error) {
	if id == "" {
		return repository.Runbook{}, fmt.Errorf("runbook id is required")
	}
	return s.repo.GetRunbook(ctx, id)
}

func (s *Service) SearchPastIncidents(ctx context.Context, query, service string, limit int) ([]repository.PastIncident, error) {
	return s.repo.SearchPastIncidents(ctx, repository.IncidentFilter{
		Query: query, Service: service, Limit: limit,
	})
}
