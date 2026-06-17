package config

import (
	"flag"
	"os"
)

const (
	ServerMonitoring    = "monitoring"
	ServerDeployments   = "deployments"
	ServerLogs          = "logs"
	ServerKnowledge     = "knowledge"
	ServerNotifications = "notifications"
)

// Config holds orchestrator runtime configuration.
type Config struct {
	Orchestrator OrchestratorConfig
	MCP          MCPEndpoints
	Log          LogConfig
}

type OrchestratorConfig struct {
	Name            string
	Version         string
	DefaultService  string
	ConnectTimeoutS int
}

type MCPEndpoints struct {
	Monitoring    string
	Deployments   string
	Logs          string
	Knowledge     string
	Notifications string
}

type LogConfig struct {
	Filename string
}

// Endpoints returns a map of server name → HTTP endpoint.
func (m MCPEndpoints) Endpoints() map[string]string {
	return map[string]string{
		ServerMonitoring:    m.Monitoring,
		ServerDeployments:   m.Deployments,
		ServerLogs:          m.Logs,
		ServerKnowledge:     m.Knowledge,
		ServerNotifications: m.Notifications,
	}
}

func Default() *Config {
	return &Config{
		Orchestrator: OrchestratorConfig{
			Name:            "incident-pilot-orchestrator",
			Version:         "0.1.0",
			DefaultService:  "payment-api",
			ConnectTimeoutS: 10,
		},
		MCP: MCPEndpoints{
			Monitoring:    "http://localhost:8081",
			Deployments:   "http://localhost:8082",
			Logs:          "http://localhost:8083",
			Knowledge:     "http://localhost:8084",
			Notifications: "http://localhost:8085",
		},
		Log: LogConfig{Filename: "/tmp/orchestrator.log"},
	}
}

func Load() *Config {
	cfg := Default()

	service := flag.String("service", envOrDefault("INCIDENT_SERVICE", cfg.Orchestrator.DefaultService), "Service under investigation")
	logFile := flag.String("log-file", envOrDefault("LOG_FILE", cfg.Log.Filename), "Log file path")
	monitoring := flag.String("mcp-monitoring", envOrDefault("MCP_MONITORING_URL", cfg.MCP.Monitoring), "Monitoring MCP URL")
	deployments := flag.String("mcp-deployments", envOrDefault("MCP_DEPLOYMENTS_URL", cfg.MCP.Deployments), "Deployments MCP URL")
	logs := flag.String("mcp-logs", envOrDefault("MCP_LOGS_URL", cfg.MCP.Logs), "Logs MCP URL")
	knowledge := flag.String("mcp-knowledge", envOrDefault("MCP_KNOWLEDGE_URL", cfg.MCP.Knowledge), "Knowledge MCP URL")
	notifications := flag.String("mcp-notifications", envOrDefault("MCP_NOTIFICATIONS_URL", cfg.MCP.Notifications), "Notifications MCP URL")
	flag.Parse()

	cfg.Orchestrator.DefaultService = *service
	cfg.Log.Filename = *logFile
	cfg.MCP.Monitoring = *monitoring
	cfg.MCP.Deployments = *deployments
	cfg.MCP.Logs = *logs
	cfg.MCP.Knowledge = *knowledge
	cfg.MCP.Notifications = *notifications

	return cfg
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
