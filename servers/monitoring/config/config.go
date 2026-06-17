package config

import (
	"flag"
	"os"
	"time"
)

const (
	TransportHTTP  = "http"
	TransportStdio = "stdio"
)

// Config holds runtime configuration for the monitoring MCP server.
type Config struct {
	Server    ServerConfig
	Transport TransportConfig
	Log       LogConfig
	Auth      AuthConfig
}

type ServerConfig struct {
	Name    string
	Version string
}

type TransportConfig struct {
	Mode            string
	HTTPAddr        string
	ShutdownTimeout time.Duration
}

type LogConfig struct {
	Filename string
}

// AuthConfig is reserved for future authorization (API keys, OAuth, RBAC).
type AuthConfig struct {
	Enabled bool
	APIKey  string
}

// Default returns production-oriented defaults.
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Name:    "incident-pilot-monitoring",
			Version: "0.1.0",
		},
		Transport: TransportConfig{
			Mode:            TransportHTTP,
			HTTPAddr:        ":8081",
			ShutdownTimeout: 5 * time.Second,
		},
		Log: LogConfig{
			Filename: "/tmp/monitoring-mcp.log",
		},
	}
}

// Load parses flags and environment variables into Config.
func Load() *Config {
	cfg := Default()

	transport := flag.String("transport", envOrDefault("MCP_TRANSPORT", cfg.Transport.Mode), "Transport: http or stdio")
	addr := flag.String("addr", envOrDefault("MCP_ADDR", cfg.Transport.HTTPAddr), "HTTP listen address")
	logFile := flag.String("log-file", envOrDefault("LOG_FILE", cfg.Log.Filename), "Log file path")
	authEnabled := flag.Bool("auth-enabled", envOrDefaultBool("MCP_AUTH_ENABLED", false), "Enable authorization (reserved)")
	authAPIKey := flag.String("auth-api-key", os.Getenv("MCP_AUTH_API_KEY"), "API key for authorization (reserved)")
	flag.Parse()

	cfg.Transport.Mode = *transport
	cfg.Transport.HTTPAddr = *addr
	cfg.Log.Filename = *logFile
	cfg.Auth.Enabled = *authEnabled
	cfg.Auth.APIKey = *authAPIKey

	return cfg
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrDefaultBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	switch v {
	case "1", "true", "TRUE", "yes", "YES":
		return true
	case "0", "false", "FALSE", "no", "NO":
		return false
	default:
		return fallback
	}
}
