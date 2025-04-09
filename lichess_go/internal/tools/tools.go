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

type LichessTools struct {
	Server *server.MCPServer
	API    *lichess.Lichess
}

func NewLichessTools(server *server.MCPServer, api *lichess.Lichess) *LichessTools {
	return &LichessTools{
		Server: server,
		API:    api,
	}
}

func (t *LichessTools) RegisterTools() {
	t.Server.AddTool(t.BoardStatus())
}
