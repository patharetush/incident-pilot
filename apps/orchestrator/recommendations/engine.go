package recommendations

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/patharetush/incident-pilot/apps/orchestrator/evidence"
)

type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)

// Recommendation is an actionable suggestion derived from collected evidence.
type Recommendation struct {
	ID               string   `json:"id"`
	Priority         Priority `json:"priority"`
	Action           string   `json:"action"`
	Rationale        string   `json:"rationale"`
	RequiresApproval bool     `json:"requires_approval"`
	EvidenceIDs      []string `json:"evidence_ids,omitempty"`
	Domain           string   `json:"domain"`
}

// Engine produces recommendations from evidence using rule-based analysis.
type Engine struct{}

func New() *Engine { return &Engine{} }

func (e *Engine) Analyze(service string, items []evidence.Item) []Recommendation {
	var recs []Recommendation
	seq := 0
	nextID := func() string {
		seq++
		return fmt.Sprintf("rec-%03d", seq)
	}

	hasKind := func(k evidence.Kind) ([]evidence.Item, bool) {
		var matched []evidence.Item
		for _, item := range items {
			if item.Kind == k {
				matched = append(matched, item)
			}
		}
		return matched, len(matched) > 0
	}

	if alerts, ok := hasKind(evidence.KindAlert); ok {
		recs = append(recs, Recommendation{
			ID: nextID(), Priority: PriorityCritical,
			Action:      "Triage critical alerts and confirm blast radius",
			Rationale:   fmt.Sprintf("Found %d alert evidence item(s) during investigation of %s", len(alerts), service),
			Domain:      "monitoring",
			EvidenceIDs: idsFromItems(alerts),
		})
	}

	if metrics, ok := hasKind(evidence.KindMetric); ok {
		if highErrorRate(metrics) {
			recs = append(recs, Recommendation{
				ID: nextID(), Priority: PriorityHigh,
				Action:      "Investigate elevated error rate; correlate with recent deployments",
				Rationale:   "Service metrics indicate degraded error rate beyond normal baseline",
				Domain:      "monitoring",
				EvidenceIDs: idsFromItems(metrics),
			})
		}
	}

	if deps, ok := hasKind(evidence.KindDeployment); ok {
		recs = append(recs, Recommendation{
			ID: nextID(), Priority: PriorityHigh,
			Action:      "Review recent deployment changes for correlation with incident timeline",
			Rationale:   fmt.Sprintf("Found %d deployment/change evidence item(s) in the lookback window", len(deps)),
			Domain:      "deployments",
			EvidenceIDs: idsFromItems(deps),
		})
	}

	if logs, ok := hasKind(evidence.KindLog); ok {
		recs = append(recs, Recommendation{
			ID: nextID(), Priority: PriorityHigh,
			Action:      "Analyze error logs and trace IDs for root cause signals",
			Rationale:   fmt.Sprintf("Collected %d log evidence item(s) with error patterns", len(logs)),
			Domain:      "logs",
			EvidenceIDs: idsFromItems(logs),
		})
	}

	if runbooks, ok := hasKind(evidence.KindRunbook); ok {
		recs = append(recs, Recommendation{
			ID: nextID(), Priority: PriorityMedium,
			Action:      "Follow matching runbook steps before attempting remediation",
			Rationale:   fmt.Sprintf("Found %d applicable runbook(s) in knowledge base", len(runbooks)),
			Domain:      "knowledge",
			EvidenceIDs: idsFromItems(runbooks),
		})
	}

	if past, ok := hasKind(evidence.KindIncident); ok {
		recs = append(recs, Recommendation{
			ID: nextID(), Priority: PriorityMedium,
			Action:      "Compare with similar past incidents and proven resolutions",
			Rationale:   fmt.Sprintf("Found %d similar historical incident(s)", len(past)),
			Domain:      "knowledge",
			EvidenceIDs: idsFromItems(past),
		})
	}

	if notifs, ok := hasKind(evidence.KindNotification); ok {
		recs = append(recs, Recommendation{
			ID: nextID(), Priority: PriorityMedium,
			Action:           "Review and approve pending stakeholder notifications",
			Rationale:        "Pending notifications require human approval before send",
			RequiresApproval: true,
			Domain:           "notifications",
			EvidenceIDs:      idsFromItems(notifs),
		})
	}

	if len(recs) == 0 {
		recs = append(recs, Recommendation{
			ID: nextID(), Priority: PriorityLow,
			Action:    "Continue monitoring; insufficient evidence for automated recommendations",
			Rationale: "Expand investigation scope or wait for additional signals",
			Domain:    "orchestrator",
		})
	}

	return recs
}

func idsFromItems(items []evidence.Item) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}

func highErrorRate(items []evidence.Item) bool {
	for _, item := range items {
		var payload map[string]any
		if json.Unmarshal(item.Data, &payload) != nil {
			continue
		}
		metrics, ok := payload["metrics"].([]any)
		if !ok {
			continue
		}
		for _, m := range metrics {
			metric, ok := m.(map[string]any)
			if !ok {
				continue
			}
			name, _ := metric["name"].(string)
			value, _ := metric["value"].(float64)
			if strings.EqualFold(name, "error_rate") && value > 5 {
				return true
			}
		}
	}
	return false
}
