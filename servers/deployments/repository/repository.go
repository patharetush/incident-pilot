package repository

import (
	"context"
	"time"
)

type Repository interface {
	ListDeployments(ctx context.Context, filter DeploymentFilter) ([]Deployment, error)
	GetDeployment(ctx context.Context, id string) (Deployment, error)
	ListRecentChanges(ctx context.Context, service string, since time.Duration) ([]Change, error)
}

type DeploymentStatus string

const (
	StatusSuccess    DeploymentStatus = "success"
	StatusFailed     DeploymentStatus = "failed"
	StatusInProgress DeploymentStatus = "in_progress"
	StatusRolledBack DeploymentStatus = "rolled_back"
)

type Deployment struct {
	ID          string           `json:"id"`
	Service     string           `json:"service"`
	Version     string           `json:"version"`
	Status      DeploymentStatus `json:"status"`
	DeployedAt  time.Time        `json:"deployed_at"`
	DeployedBy  string           `json:"deployed_by"`
	Commit      string           `json:"commit"`
	Summary     string           `json:"summary"`
	Environment string           `json:"environment"`
	Changes     []Change         `json:"changes,omitempty"`
}

type Change struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Component   string `json:"component,omitempty"`
}

type DeploymentFilter struct {
	Service string
	Limit   int
}
