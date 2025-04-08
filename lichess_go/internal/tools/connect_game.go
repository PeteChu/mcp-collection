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
	t.Server.AddTool(t.ConnectGame())
}

func (l *LichessTools) ConnectGame() (mcp.Tool, MCPHandler) {
	tool := mcp.NewTool("connect_game",
		mcp.WithDescription("connect to a game with gameId"),
		mcp.WithString(
			"gameId",
			mcp.Required(),
			mcp.Description("gameId of the game to connect to"),
		),
	)
	return tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		gameId := request.Params.Arguments["gameId"].(string)
		status := l.API.GetBoard(gameId)
		return mcp.NewToolResultText(status), nil
	}
}
