# Knowledge MCP Server

MCP server for runbooks, past incidents, and mitigation knowledge during incident response.

**Default port:** `:8084`

## MCP surface

| Type | Name | Description |
|------|------|-------------|
| Tool | `search_runbooks` | Search runbooks by keyword and service |
| Tool | `get_runbook` | Full runbook with remediation steps |
| Tool | `search_past_incidents` | Historical incidents and resolutions |
| Resource | `knowledge://runbooks/catalog` | JSON catalog of runbooks |
| Resource | `knowledge://runbooks/{id}` | Single runbook detail |
| Prompt | `recommend_mitigation` | Mitigation recommendations from runbooks + history |

## Package layout

```text
servers/knowledge/
├── app.go
├── app_test.go
├── cmd/knowledge/main.go
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
| `repository` | Runbook and incident knowledge store |
| `service` | Search logic and validation |
| `tools` | MCP tool handlers |
| `resources` | MCP resources for runbook catalogs |
| `prompts` | Guided mitigation recommendation flows |
| `shared/transport` | HTTP/stdio serving and auth middleware hook |

## Run locally

```bash
make run-knowledge
# or
go run ./servers/knowledge/cmd/knowledge
```

## Tests

```bash
go test ./servers/knowledge/... -v
```

## Extension points

- Back repository with Confluence, Notion, or an internal wiki API.
- Add vector search for semantic runbook matching.
- Version runbooks and track which version was used during an incident (audit trail).
