BINARY = incident-pilot

.PHONY: all test tidy build clean \
	run-orchestrator \
	run-monitoring run-deployments run-logs run-knowledge run-notifications \
	build-monitoring build-deployments build-logs build-knowledge build-notifications

all: test

build:
	go build ./...

test:
	go test ./...

tidy:
	go mod tidy

run-orchestrator:
	@echo "Starting orchestrator..."
	go run ./apps/orchestrator

build-monitoring:
	@echo "Building monitoring MCP service..."
	go build -o monitoring-server ./servers/monitoring/cmd/monitoring

run-monitoring:
	@echo "Starting monitoring MCP service on :8081..."
	go run ./servers/monitoring/cmd/monitoring

build-deployments:
	@echo "Building deployments MCP service..."
	go build -o deployments-server ./servers/deployments/cmd/deployments

run-deployments:
	@echo "Starting deployments MCP service on :8082..."
	go run ./servers/deployments/cmd/deployments

build-logs:
	@echo "Building logs MCP service..."
	go build -o logs-server ./servers/logs/cmd/logs

run-logs:
	@echo "Starting logs MCP service on :8083..."
	go run ./servers/logs/cmd/logs

build-knowledge:
	@echo "Building knowledge MCP service..."
	go build -o knowledge-server ./servers/knowledge/cmd/knowledge

run-knowledge:
	@echo "Starting knowledge MCP service on :8084..."
	go run ./servers/knowledge/cmd/knowledge

build-notifications:
	@echo "Building notifications MCP service..."
	go build -o notifications-server ./servers/notifications/cmd/notifications

run-notifications:
	@echo "Starting notifications MCP service on :8085..."
	go run ./servers/notifications/cmd/notifications

clean:
	@echo "Resetting workspace to a clean state..."
	rm -f monitoring-server deployments-server logs-server knowledge-server notifications-server
