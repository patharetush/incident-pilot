# Logs MCP Server

MCP server for log search, trace correlation, and error pattern analysis during incidents.

**Default port:** `:8083`

## MCP surface

| Type | Name | Description |
|------|------|-------------|
| Tool | `search_logs` | Search by service, level, and message query |
| Tool | `get_log_entry` | Fetch a single log entry by ID |
| Tool | `get_log_context` | Correlated entries by trace ID |
| Tool | `list_error_patterns` | Recurring error patterns in recent logs |
| Resource | `logs://patterns/errors` | JSON summary of error patterns |
| Resource | `logs://service/{name}/recent` | Recent logs for a service |
| Prompt | `analyze_logs` | Structured log analysis workflow |

## Package layout

```text
servers/logs/
├── app.go
├── app_test.go
├── cmd/logs/main.go
├── config/
├── prompts/
├── repository/ (+ memory/)
├── resources/
├── service/
└── tools/
```

Shared runtime: `shared/config`, `shared/transport`, `shared/logging`

## Layer responsibilities

| Layer | Role |
|-------|------|
| `config` | Server-specific defaults (name, port, log file) |
| `repository` | Log entry and pattern storage interface |
| `service` | Query validation, trace correlation logic |
| `tools` | MCP tool handlers |
| `resources` | MCP resources for log snapshots |
| `prompts` | Guided log analysis workflows |
| `shared/transport` | HTTP/stdio serving and auth middleware hook |

## Run locally

```bash
make run-logs
# or
go run ./servers/logs/cmd/logs
```

## Tests

```bash
go test ./servers/logs/... -v
```

## Extension points

- Replace in-memory backend with Elasticsearch, Loki, CloudWatch Logs, or Datadog.
- Add rate limiting and query cost guards in `service`.
- Extend `get_log_context` with time-window based correlation.
