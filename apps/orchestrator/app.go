package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/patharetush/incident-pilot/apps/orchestrator/config"
	"github.com/patharetush/incident-pilot/apps/orchestrator/discovery"
	"github.com/patharetush/incident-pilot/apps/orchestrator/evidence"
	"github.com/patharetush/incident-pilot/apps/orchestrator/executor"
	"github.com/patharetush/incident-pilot/apps/orchestrator/planner"
	"github.com/patharetush/incident-pilot/apps/orchestrator/recommendations"
	"github.com/patharetush/incident-pilot/apps/orchestrator/sessions"
	"github.com/patharetush/incident-pilot/shared/logging"
)

// App is the composition root for the incident investigation orchestrator.
type App struct {
	cfg          *config.Config
	discovery    *discovery.Pool
	sessions     *sessions.Store
	planner      *planner.Planner
	executor     *executor.Executor
	evidence     *evidence.Collector
	recommender  *recommendations.Engine
}

// New wires all orchestrator layers.
func New(cfg *config.Config) *App {
	pool := discovery.NewPool()
	return &App{
		cfg:         cfg,
		discovery:   pool,
		sessions:    sessions.NewStore(),
		planner:     planner.New(),
		executor:    executor.New(pool),
		evidence:    evidence.NewCollector(),
		recommender: recommendations.New(),
	}
}

// Connect discovers and connects to all configured MCP servers.
func (a *App) Connect(ctx context.Context) error {
	return a.discovery.Connect(ctx, a.cfg.MCP.Endpoints())
}

// Close releases MCP connections.
func (a *App) Close() {
	a.discovery.Close()
}

// Catalog returns discovered MCP server capabilities.
func (a *App) Catalog() []*discovery.ServerCatalog {
	return a.discovery.Catalog()
}

// GetSession retrieves a session by ID.
func (a *App) GetSession(id string) (*sessions.Session, error) {
	return a.sessions.Get(id)
}

// RunInvestigation executes the full investigate → evidence → recommend pipeline.
func (a *App) RunInvestigation(ctx context.Context, service string) (*sessions.Session, error) {
	if service == "" {
		return nil, fmt.Errorf("service is required")
	}

	session := a.sessions.Create(service)
	session.Status = sessions.StatusPlanning
	a.sessions.Update(session)

	logging.L().Info().Str("session", session.ID).Str("service", service).Msg("starting investigation")

	plan := a.planner.Build(service)
	session.Plan = plan
	session.Status = sessions.StatusExecuting
	a.sessions.Update(session)

	outcomes := a.executor.Run(ctx, plan)

	session.Status = sessions.StatusAnalyzing
	var allEvidence []evidence.Item
	for _, outcome := range outcomes {
		items := a.evidence.Collect(session.ID, outcome.StepID, outcome.Server, outcome.Tool, outcome.Result)
		allEvidence = append(allEvidence, items...)
	}
	session.Evidence = allEvidence

	recs := a.recommender.Analyze(service, allEvidence)
	session.Recommendations = recs

	now := time.Now().UTC().Truncate(time.Second)
	session.Status = sessions.StatusCompleted
	session.CompletedAt = &now
	a.sessions.Update(session)

	logging.L().Info().
		Str("session", session.ID).
		Int("evidence", len(allEvidence)).
		Int("recommendations", len(recs)).
		Msg("investigation completed")

	return session, nil
}

// PrintReport writes a human-readable investigation report to stdout.
func PrintReport(session *sessions.Session) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(session)
}
