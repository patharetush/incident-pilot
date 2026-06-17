package planner

import (
	"fmt"
	"time"

	"github.com/patharetush/incident-pilot/apps/orchestrator/config"
)

type StepStatus string

const (
	StepPending   StepStatus = "pending"
	StepRunning   StepStatus = "running"
	StepCompleted StepStatus = "completed"
	StepFailed    StepStatus = "failed"
	StepSkipped   StepStatus = "skipped"
)

// Step is a single investigation action against an MCP server tool.
type Step struct {
	ID          string         `json:"id"`
	Server      string         `json:"server"`
	Tool        string         `json:"tool"`
	Arguments   map[string]any `json:"arguments,omitempty"`
	Description string         `json:"description"`
	Status      StepStatus     `json:"status"`
	Error       string         `json:"error,omitempty"`
}

// Plan is an ordered investigation workflow.
type Plan struct {
	ID        string    `json:"id"`
	Service   string    `json:"service"`
	Steps     []Step    `json:"steps"`
	CreatedAt time.Time `json:"created_at"`
}

// Planner builds rule-based investigation plans. LLM-backed planning can replace this layer.
type Planner struct{}

func New() *Planner { return &Planner{} }

// Build creates a standard multi-domain investigation plan for a service.
func (p *Planner) Build(service string) *Plan {
	now := time.Now().UTC().Truncate(time.Second)
	steps := []Step{
		{
			ID: "step-01", Server: config.ServerMonitoring, Tool: "list_alerts",
			Arguments: map[string]any{"severity": "critical"},
			Description: "Identify critical alerts across the fleet",
		},
		{
			ID: "step-02", Server: config.ServerMonitoring, Tool: "get_service_metrics",
			Arguments: map[string]any{"service": service},
			Description: fmt.Sprintf("Collect live metrics for %s", service),
		},
		{
			ID: "step-03", Server: config.ServerDeployments, Tool: "list_recent_changes",
			Arguments: map[string]any{"service": service, "since_minutes": 120},
			Description: fmt.Sprintf("Find recent deployments and config changes for %s", service),
		},
		{
			ID: "step-04", Server: config.ServerLogs, Tool: "search_logs",
			Arguments: map[string]any{"service": service, "level": "error", "limit": 10},
			Description: fmt.Sprintf("Search error logs for %s", service),
		},
		{
			ID: "step-05", Server: config.ServerLogs, Tool: "list_error_patterns",
			Arguments: map[string]any{"service": service},
			Description: fmt.Sprintf("Detect recurring error patterns for %s", service),
		},
		{
			ID: "step-06", Server: config.ServerKnowledge, Tool: "search_runbooks",
			Arguments: map[string]any{"service": service, "query": "error"},
			Description: fmt.Sprintf("Find applicable runbooks for %s", service),
		},
		{
			ID: "step-07", Server: config.ServerKnowledge, Tool: "search_past_incidents",
			Arguments: map[string]any{"service": service, "query": "error"},
			Description: fmt.Sprintf("Search similar past incidents for %s", service),
		},
		{
			ID: "step-08", Server: config.ServerNotifications, Tool: "list_pending_approvals",
			Description: "Check pending stakeholder notifications awaiting approval",
		},
	}
	return &Plan{
		ID: fmt.Sprintf("plan-%s-%d", service, now.Unix()),
		Service: service,
		Steps:   steps,
		CreatedAt: now,
	}
}
