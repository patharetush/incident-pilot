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

// Config holds runtime configuration shared by all MCP servers.
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

// Defaults identifies server-specific baseline settings.
type Defaults struct {
	ServerName    string
	ServerVersion string
	HTTPAddr      string
	LogFilename   string
}

// Default returns production-oriented defaults for a server.
func Default(d Defaults) *Config {
	version := d.ServerVersion
	if version == "" {
		version = "0.1.0"
	}
	return &Config{
		Server: ServerConfig{
			Name:    d.ServerName,
			Version: version,
		},
		Transport: TransportConfig{
			Mode:            TransportHTTP,
			HTTPAddr:        d.HTTPAddr,
			ShutdownTimeout: 5 * time.Second,
		},
		Log: LogConfig{
			Filename: d.LogFilename,
		},
	}
}

// Load parses flags and environment variables into Config.
func Load(d Defaults) *Config {
	cfg := Default(d)

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
