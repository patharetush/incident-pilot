package repository

import (
	"context"
	"time"
)

type Repository interface {
	SearchRunbooks(ctx context.Context, filter RunbookFilter) ([]Runbook, error)
	GetRunbook(ctx context.Context, id string) (Runbook, error)
	SearchPastIncidents(ctx context.Context, filter IncidentFilter) ([]PastIncident, error)
}

type Runbook struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Service   string   `json:"service"`
	Tags      []string `json:"tags"`
	Symptoms  []string `json:"symptoms"`
	Steps     []string `json:"steps"`
	Summary   string   `json:"summary"`
}

type PastIncident struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Service    string    `json:"service"`
	Summary    string    `json:"summary"`
	Resolution string    `json:"resolution"`
	OccurredAt time.Time `json:"occurred_at"`
	Tags       []string  `json:"tags"`
}

type RunbookFilter struct {
	Query   string
	Service string
	Limit   int
}

type IncidentFilter struct {
	Query   string
	Service string
	Limit   int
}
