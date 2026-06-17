package service

import (
	"context"
	"fmt"
	"time"

	"github.com/patharetush/incident-pilot/servers/deployments/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListDeployments(ctx context.Context, service string, limit int) ([]repository.Deployment, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListDeployments(ctx, repository.DeploymentFilter{
		Service: service,
		Limit:   limit,
	})
}

func (s *Service) GetDeployment(ctx context.Context, id string) (repository.Deployment, error) {
	if id == "" {
		return repository.Deployment{}, fmt.Errorf("deployment id is required")
	}
	return s.repo.GetDeployment(ctx, id)
}

func (s *Service) ListRecentChanges(ctx context.Context, service string, sinceMinutes int) ([]repository.Change, error) {
	if service == "" {
		return nil, fmt.Errorf("service is required")
	}
	if sinceMinutes <= 0 {
		sinceMinutes = 60
	}
	return s.repo.ListRecentChanges(ctx, service, time.Duration(sinceMinutes)*time.Minute)
}
