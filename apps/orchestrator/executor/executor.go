package executor

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/patharetush/incident-pilot/apps/orchestrator/discovery"
	"github.com/patharetush/incident-pilot/apps/orchestrator/planner"
	"github.com/patharetush/incident-pilot/shared/logging"
)

// Outcome captures the result of a single plan step execution.
type Outcome struct {
	StepID string
	Server string
	Tool   string
	Result *mcp.CallToolResult
}

// Executor runs investigation plan steps against MCP servers.
type Executor struct {
	pool *discovery.Pool
}

func New(pool *discovery.Pool) *Executor {
	return &Executor{pool: pool}
}

// Run executes all plan steps and returns outcomes for evidence collection.
func (e *Executor) Run(ctx context.Context, plan *planner.Plan) []Outcome {
	var outcomes []Outcome

	for i := range plan.Steps {
		step := &plan.Steps[i]
		if !e.pool.HasServer(step.Server) {
			step.Status = planner.StepSkipped
			step.Error = fmt.Sprintf("server %q not available", step.Server)
			logging.L().Warn().Str("step", step.ID).Str("server", step.Server).Msg("step skipped")
			continue
		}

		step.Status = planner.StepRunning
		logging.L().Info().
			Str("step", step.ID).
			Str("server", step.Server).
			Str("tool", step.Tool).
			Msg("executing plan step")

		result, err := e.pool.CallTool(ctx, step.Server, step.Tool, step.Arguments)
		if err != nil {
			step.Status = planner.StepFailed
			step.Error = err.Error()
			logging.L().Error().Err(err).Str("step", step.ID).Msg("step failed")
			continue
		}
		if result.IsError {
			step.Status = planner.StepFailed
			step.Error = toolErrorMessage(result)
			continue
		}

		step.Status = planner.StepCompleted
		outcomes = append(outcomes, Outcome{
			StepID: step.ID,
			Server: step.Server,
			Tool:   step.Tool,
			Result: result,
		})
	}
	return outcomes
}

func toolErrorMessage(result *mcp.CallToolResult) string {
	for _, c := range result.Content {
		if text, ok := c.(*mcp.TextContent); ok {
			return text.Text
		}
	}
	return "tool returned error"
}
