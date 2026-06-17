package service

import (
	"context"
	"fmt"

	"github.com/patharetush/incident-pilot/servers/monitoring/repository"
)

// Service contains monitoring business logic. Authorization and policy checks
// belong here so all entry points (tools, resources, prompts) share rules.
type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListServices(ctx context.Context) ([]repository.Service, error) {
	return s.repo.ListServices(ctx)
}

func (s *Service) GetService(ctx context.Context, name string) (repository.Service, error) {
	if name == "" {
		return repository.Service{}, fmt.Errorf("service name is required")
	}
	return s.repo.GetService(ctx, name)
}

func (s *Service) GetServiceMetrics(ctx context.Context, name string) (string, []repository.MetricPoint, error) {
	if name == "" {
		return "", nil, fmt.Errorf("service name is required")
	}
	metrics, err := s.repo.GetMetrics(ctx, name)
	if err != nil {
		return "", nil, err
	}
	return name, metrics, nil
}

func (s *Service) ListAlerts(ctx context.Context, severity string) ([]repository.Alert, error) {
	return s.repo.ListAlerts(ctx, repository.AlertFilter{Severity: severity})
}

func (s *Service) GetAlert(ctx context.Context, id string) (repository.Alert, error) {
	if id == "" {
		return repository.Alert{}, fmt.Errorf("alert id is required")
	}
	return s.repo.GetAlert(ctx, id)
}
