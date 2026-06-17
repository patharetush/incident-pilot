package memory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/patharetush/incident-pilot/servers/knowledge/repository"
)

type Repository struct {
	runbooks  map[string]repository.Runbook
	incidents []repository.PastIncident
}

func New() *Repository {
	runbooks := map[string]repository.Runbook{
		"rb-001": {
			ID: "rb-001", Title: "Payment API High Error Rate",
			Service: "payment-api", Tags: []string{"5xx", "latency", "payments"},
			Symptoms: []string{"Elevated 5xx on /v1/charge", "P99 latency above SLO"},
			Steps: []string{
				"Check recent deployments for payment-api",
				"Inspect upstream payment-gateway health and timeouts",
				"Review error logs for retry exhaustion patterns",
				"Consider rollback if error rate persists >15 minutes",
			},
			Summary: "Playbook for payment-api error rate and latency incidents",
		},
		"rb-002": {
			ID: "rb-002", Title: "Database Connection Pool Exhaustion",
			Service: "inventory-db", Tags: []string{"database", "connections", "timeouts"},
			Symptoms: []string{"Connection pool at capacity", "Query timeouts increasing"},
			Steps: []string{
				"Verify active connection count vs max pool size",
				"Check for long-running queries or lock contention",
				"Scale read replicas or increase pool size if safe",
				"Identify and kill runaway sessions if necessary",
			},
			Summary: "Playbook for database connection pool incidents",
		},
		"rb-003": {
			ID: "rb-003", Title: "Service Rollback Procedure",
			Service: "", Tags: []string{"rollback", "deployment"},
			Steps: []string{
				"Confirm incident correlation with latest deployment",
				"Verify rollback target version is healthy in staging",
				"Execute rollback via deployment pipeline",
				"Monitor error rate and latency for 15 minutes post-rollback",
			},
			Summary: "Generic rollback procedure for production services",
		},
	}
	incidents := []repository.PastIncident{
		{
			ID: "inc-2024-089", Title: "Payment gateway timeout cascade",
			Service: "payment-api",
			Summary: "Upstream gateway latency caused retry storms and 5xx spike",
			Resolution: "Reduced timeout, disabled aggressive retries, rolled back v2.12 config change",
			OccurredAt: time.Now().UTC().Add(-45 * 24 * time.Hour),
			Tags: []string{"timeout", "retry", "payment-gateway"},
		},
		{
			ID: "inc-2024-112", Title: "Inventory DB pool saturation",
			Service: "inventory-db",
			Summary: "Migration left long-running lock; connection pool exhausted",
			Resolution: "Killed blocking session, increased pool size temporarily, fixed migration",
			OccurredAt: time.Now().UTC().Add(-20 * 24 * time.Hour),
			Tags: []string{"database", "migration", "connections"},
		},
	}
	return &Repository{runbooks: runbooks, incidents: incidents}
}

func (r *Repository) SearchRunbooks(ctx context.Context, filter repository.RunbookFilter) ([]repository.Runbook, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	var out []repository.Runbook
	for _, rb := range r.runbooks {
		if filter.Service != "" && rb.Service != "" && !strings.EqualFold(rb.Service, filter.Service) {
			continue
		}
		if filter.Query != "" && !matchesQuery(filter.Query, rb.Title, rb.Summary, strings.Join(rb.Tags, " ")) {
			continue
		}
		out = append(out, rb)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (r *Repository) GetRunbook(ctx context.Context, id string) (repository.Runbook, error) {
	rb, ok := r.runbooks[strings.ToLower(id)]
	if !ok {
		return repository.Runbook{}, fmt.Errorf("runbook %q not found", id)
	}
	return rb, nil
}

func (r *Repository) SearchPastIncidents(ctx context.Context, filter repository.IncidentFilter) ([]repository.PastIncident, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	var out []repository.PastIncident
	for _, inc := range r.incidents {
		if filter.Service != "" && !strings.EqualFold(inc.Service, filter.Service) {
			continue
		}
		if filter.Query != "" && !matchesQuery(filter.Query, inc.Title, inc.Summary, inc.Resolution, strings.Join(inc.Tags, " ")) {
			continue
		}
		out = append(out, inc)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func matchesQuery(query string, fields ...string) bool {
	q := strings.ToLower(query)
	for _, f := range fields {
		if strings.Contains(strings.ToLower(f), q) {
			return true
		}
	}
	return false
}
