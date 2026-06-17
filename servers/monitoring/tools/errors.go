package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func toolError(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: message},
		},
	}
}
