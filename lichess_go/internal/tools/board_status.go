package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (l *LichessTools) BoardStatus() (mcp.Tool, MCPHandler) {
	tool := mcp.NewTool("board_status",
		mcp.WithDescription("get the status of a game by gameId"),
		mcp.WithString(
			"gameId",
			mcp.Required(),
			mcp.Description("id of the game to get the status of"),
		),
	)
	return tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		gameId := request.Params.Arguments["gameId"].(string)
		status := l.API.GetBoard(gameId)
		return mcp.NewToolResultText(status), nil
	}
}
