package memory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/patharetush/incident-pilot/servers/deployments/repository"
)

type Repository struct {
	deployments map[string]repository.Deployment
}

func New() *Repository {
	now := time.Now().UTC().Truncate(time.Second)
	deployments := map[string]repository.Deployment{
		"dep-001": {
			ID:          "dep-001",
			Service:     "payment-api",
			Version:     "v2.14.0",
			Status:      repository.StatusSuccess,
			DeployedAt:  now.Add(-18 * time.Minute),
			DeployedBy:  "ci-pipeline",
			Commit:      "a1b2c3d",
			Summary:     "Payment retry logic update and timeout tuning",
			Environment: "production",
			Changes: []repository.Change{
				{Type: "config", Description: "Increased upstream timeout from 500ms to 800ms", Component: "payment-gateway-client"},
				{Type: "code", Description: "Added retry on 502/503 responses", Component: "payment-handler"},
			},
		},
		"dep-002": {
			ID:          "dep-002",
			Service:     "payment-api",
			Version:     "v2.13.2",
			Status:      repository.StatusRolledBack,
			DeployedAt:  now.Add(-26 * time.Hour),
			DeployedBy:  "ci-pipeline",
			Commit:      "f9e8d7c",
			Summary:     "Hotfix for currency rounding edge case",
			Environment: "production",
		},
		"dep-003": {
			ID:          "dep-003",
			Service:     "user-service",
			Version:     "v1.8.1",
			Status:      repository.StatusSuccess,
			DeployedAt:  now.Add(-3 * time.Hour),
			DeployedBy:  "ci-pipeline",
			Commit:      "b4c5d6e",
			Summary:     "Profile cache TTL adjustment",
			Environment: "production",
		},
		"dep-004": {
			ID:          "dep-004",
			Service:     "inventory-db",
			Version:     "migration-042",
			Status:      repository.StatusFailed,
			DeployedAt:  now.Add(-45 * time.Minute),
			DeployedBy:  "db-migrations",
			Commit:      "m042-index",
			Summary:     "Index rebuild migration — rolled back automatically",
			Environment: "production",
			Changes: []repository.Change{
				{Type: "migration", Description: "Rebuild index on inventory_items(sku)", Component: "inventory-db"},
			},
		},
	}
	return &Repository{deployments: deployments}
}

func (r *Repository) ListDeployments(ctx context.Context, filter repository.DeploymentFilter) ([]repository.Deployment, error) {
	out := make([]repository.Deployment, 0)
	for _, dep := range r.deployments {
		if filter.Service != "" && !strings.EqualFold(dep.Service, filter.Service) {
			continue
		}
		out = append(out, dep)
	}
	if filter.Limit > 0 && len(out) > filter.Limit {
		out = out[:filter.Limit]
	}
	return out, nil
}

func (r *Repository) GetDeployment(ctx context.Context, id string) (repository.Deployment, error) {
	dep, ok := r.deployments[strings.ToLower(id)]
	if !ok {
		return repository.Deployment{}, fmt.Errorf("deployment %q not found", id)
	}
	return dep, nil
}

func (r *Repository) ListRecentChanges(ctx context.Context, service string, since time.Duration) ([]repository.Change, error) {
	if service == "" {
		return nil, fmt.Errorf("service is required")
	}
	cutoff := time.Now().UTC().Add(-since)
	var out []repository.Change
	for _, dep := range r.deployments {
		if !strings.EqualFold(dep.Service, service) {
			continue
		}
		if dep.DeployedAt.Before(cutoff) {
			continue
		}
		for _, ch := range dep.Changes {
			out = append(out, ch)
		}
	}
	if len(out) == 0 {
		return []repository.Change{}, nil
	}
	return out, nil
}
