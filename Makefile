BINARY = incident-pilot

.PHONY: all test tidy run-orchestrator run-monitoring run-logs run-deployments run-knowledge run-notifications build

all: test

build:
	go build ./...

test:
	go test ./...

tidy:
	go mod tidy

build-monitoring:
	@echo "Building monitoring MCP service..."
	go build -o monitoring-server ./servers/monitoring/cmd/monitoring

run-monitoring:
	@echo "Starting monitoring MCP service..."
	go run ./servers/monitoring/cmd/monitoring

clean:
	@echo "Resetting workspace to a clean state..."
	rm -rf monitoring-server
