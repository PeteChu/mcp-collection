package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (t *MCPTools) BoardMove() (mcp.Tool, MCPHandler) {
	tool := mcp.NewTool("board_move",
		mcp.WithDescription("make a move on the board"),
		mcp.WithString(
			"gameId",
			mcp.Required(),
			mcp.Description("id of the game to make a move on"),
		),
		mcp.WithString(
			"move",
			mcp.Required(),
			mcp.Description("The move to play, in UCI format"),
		),
		mcp.WithBoolean("offerDraw",
			mcp.Description("Whether to offer (or agree to) a draw"),
		),
	)
	return tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		gameId := request.Params.Arguments["gameId"].(string)
		move := request.Params.Arguments["move"].(string)

		draw := false
		if _, ok := request.Params.Arguments["offerDraw"]; ok {
			draw = true
		}

		status := t.API.MakeMove(gameId, move, draw)
		return mcp.NewToolResultText(status), nil
	}
}
