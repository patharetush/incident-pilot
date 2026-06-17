# Deployments MCP Server

MCP server for deployment history, change correlation, and rollback evidence during incident response.

**Default port:** `:8082`

## MCP surface

| Type | Name | Description |
|------|------|-------------|
| Tool | `list_deployments` | Recent deployments, optional service filter |
| Tool | `get_deployment` | Full deployment record by ID |
| Tool | `list_recent_changes` | Config/code/migration changes in a time window |
| Resource | `deployments://catalog/recent` | JSON catalog of recent deployments |
| Resource | `deployments://deployments/{id}` | Single deployment detail |
| Prompt | `correlate_deployment` | Correlate deployments with an active incident |

## Package layout

```text
servers/deployments/
├── app.go
├── app_test.go
├── cmd/deployments/main.go
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
| `repository` | Deployment data access interface + in-memory backend |
| `service` | Validation and business logic |
| `tools` | MCP tool handlers (thin adapters) |
| `resources` | MCP resources for structured snapshots |
| `prompts` | Guided deployment correlation workflows |
| `shared/transport` | HTTP/stdio serving and auth middleware hook |

## Run locally

```bash
make run-deployments
# or
go run ./servers/deployments/cmd/deployments
```

## Tests

```bash
go test ./servers/deployments/... -v
```

## Extension points

- Swap `repository/memory` for CI/CD integrations (GitHub Actions, Argo CD, Spinnaker).
- Inject custom repository via `deployments.New(cfg, &deployments.Options{Repository: repo})`.
- Add approval gates in `service` before exposing destructive deployment actions.
