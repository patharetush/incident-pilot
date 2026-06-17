package mcptest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ConnectHTTP starts an in-process HTTP MCP session for end-to-end tests.
func ConnectHTTP(t *testing.T, server *mcp.Server) *mcp.ClientSession {
	t.Helper()

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{JSONResponse: true})

	httpServer := httptest.NewServer(handler)
	t.Cleanup(httpServer.Close)

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "incident-pilot-test-client",
		Version: "0.1.0",
	}, nil)

	session, err := client.Connect(t.Context(), &mcp.StreamableClientTransport{
		Endpoint: httpServer.URL,
	}, nil)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

// DecodeStructuredOutput unmarshals MCP structured tool output into dest.
func DecodeStructuredOutput(t *testing.T, result *mcp.CallToolResult, dest any) {
	t.Helper()
	if result.StructuredContent == nil {
		t.Fatal("expected structured content in tool result")
	}
	data, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		t.Fatalf("unmarshal structured content: %v", err)
	}
}
