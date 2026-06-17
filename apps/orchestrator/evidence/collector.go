package evidence

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Kind string

const (
	KindAlert        Kind = "alert"
	KindMetric       Kind = "metric"
	KindLog          Kind = "log"
	KindDeployment   Kind = "deployment"
	KindRunbook      Kind = "runbook"
	KindIncident     Kind = "incident"
	KindNotification Kind = "notification"
	KindGeneric      Kind = "generic"
)

// Item is a structured evidence record collected during investigation.
type Item struct {
	ID          string          `json:"id"`
	SessionID   string          `json:"session_id"`
	StepID      string          `json:"step_id"`
	Source      string          `json:"source"`
	Tool        string          `json:"tool"`
	Kind        Kind            `json:"kind"`
	Summary     string          `json:"summary"`
	Data        json.RawMessage `json:"data,omitempty"`
	CollectedAt time.Time       `json:"collected_at"`
}

// Collector extracts evidence from MCP tool results.
type Collector struct {
	seq int
}

func NewCollector() *Collector { return &Collector{} }

func (c *Collector) Collect(sessionID string, stepID, server, tool string, result *mcp.CallToolResult) []Item {
	if result == nil || result.IsError {
		return nil
	}
	now := time.Now().UTC().Truncate(time.Second)
	kind := classifyTool(tool)

	if result.StructuredContent != nil {
		data, err := json.Marshal(result.StructuredContent)
		if err != nil {
			return nil
		}
		c.seq++
		return []Item{{
			ID: c.id(), SessionID: sessionID, StepID: stepID,
			Source: server, Tool: tool, Kind: kind,
			Summary: summarizeStructured(server, tool, result.StructuredContent),
			Data:    data, CollectedAt: now,
		}}
	}

	var items []Item
	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok && text.Text != "" {
			c.seq++
			raw, _ := json.Marshal(map[string]string{"text": text.Text})
			items = append(items, Item{
				ID: c.id(), SessionID: sessionID, StepID: stepID,
				Source: server, Tool: tool, Kind: kind,
				Summary: truncate(text.Text, 200),
				Data:    raw, CollectedAt: now,
			})
		}
	}
	return items
}

func (c *Collector) id() string {
	return fmt.Sprintf("evd-%04d", c.seq)
}

func classifyTool(tool string) Kind {
	switch tool {
	case "list_alerts", "get_alert":
		return KindAlert
	case "get_service_metrics", "list_services":
		return KindMetric
	case "search_logs", "get_log_entry", "get_log_context", "list_error_patterns":
		return KindLog
	case "list_deployments", "get_deployment", "list_recent_changes":
		return KindDeployment
	case "search_runbooks", "get_runbook":
		return KindRunbook
	case "search_past_incidents":
		return KindIncident
	case "list_notifications", "list_pending_approvals", "list_channels":
		return KindNotification
	default:
		return KindGeneric
	}
}

func summarizeStructured(server, tool string, content any) string {
	data, err := json.Marshal(content)
	if err != nil {
		return fmt.Sprintf("%s/%s result", server, tool)
	}
	var m map[string]any
	if json.Unmarshal(data, &m) == nil {
		if count, ok := m["count"].(float64); ok {
			return fmt.Sprintf("%s/%s returned %d items", server, tool, int(count))
		}
	}
	return truncate(string(data), 200)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
