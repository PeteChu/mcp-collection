package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/petechu/lichess-mcp/internal/lichess"
)

type (
	MCPHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
)

type MCPTools struct {
	Server *server.MCPServer
	API    *lichess.Lichess
}

func NewLichessTools(server *server.MCPServer, api *lichess.Lichess) *MCPTools {
	return &MCPTools{
		Server: server,
		API:    api,
	}
}

func (t *MCPTools) RegisterTools() {
	t.Server.AddTool(t.BoardStatus())
	t.Server.AddTool(t.BoardMove())
}
