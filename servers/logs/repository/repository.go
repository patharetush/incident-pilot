package repository

import (
	"context"
	"time"
)

type Repository interface {
	SearchLogs(ctx context.Context, filter LogFilter) ([]LogEntry, error)
	GetLogEntry(ctx context.Context, id string) (LogEntry, error)
	ListErrorPatterns(ctx context.Context, service string) ([]ErrorPattern, error)
}

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type LogEntry struct {
	ID        string            `json:"id"`
	Service   string            `json:"service"`
	Level     LogLevel          `json:"level"`
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	TraceID   string            `json:"trace_id,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type ErrorPattern struct {
	Pattern string `json:"pattern"`
	Count   int    `json:"count"`
	Service string `json:"service"`
	Sample  string `json:"sample"`
}

type LogFilter struct {
	Service string
	Level   LogLevel
	Query   string
	Limit   int
}
