package sessions

import (
	"fmt"
	"sync"
	"time"

	"github.com/patharetush/incident-pilot/apps/orchestrator/evidence"
	"github.com/patharetush/incident-pilot/apps/orchestrator/planner"
	"github.com/patharetush/incident-pilot/apps/orchestrator/recommendations"
)

type Status string

const (
	StatusCreated    Status = "created"
	StatusPlanning   Status = "planning"
	StatusExecuting  Status = "executing"
	StatusAnalyzing  Status = "analyzing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// Session represents an incident investigation lifecycle.
type Session struct {
	ID              string                        `json:"id"`
	Service         string                        `json:"service"`
	Status          Status                        `json:"status"`
	Plan            *planner.Plan                   `json:"plan,omitempty"`
	Evidence        []evidence.Item               `json:"evidence,omitempty"`
	Recommendations []recommendations.Recommendation `json:"recommendations,omitempty"`
	Error           string                        `json:"error,omitempty"`
	CreatedAt       time.Time                     `json:"created_at"`
	UpdatedAt       time.Time                     `json:"updated_at"`
	CompletedAt     *time.Time                    `json:"completed_at,omitempty"`
}

// Store is an in-memory session repository.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	seq      int
}

func NewStore() *Store {
	return &Store{sessions: make(map[string]*Session)}
}

func (s *Store) Create(service string) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	now := time.Now().UTC().Truncate(time.Second)
	session := &Session{
		ID:        fmt.Sprintf("sess-%04d", s.seq),
		Service:   service,
		Status:    StatusCreated,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.sessions[session.ID] = session
	return session
}

func (s *Store) Get(id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session %q not found", id)
	}
	return session, nil
}

func (s *Store) Update(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session.UpdatedAt = time.Now().UTC().Truncate(time.Second)
	s.sessions[session.ID] = session
}

func (s *Store) List() []*Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		out = append(out, session)
	}
	return out
}
