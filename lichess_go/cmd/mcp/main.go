package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/petechu/lichess-mcp/internal/lichess"
	"github.com/petechu/lichess-mcp/internal/tools"
)

func main() {
	s := server.NewMCPServer("Lichess", "0.0.1",
		server.WithToolCapabilities(true),
	)

	api := lichess.NewLichess(os.Getenv("LICHESS_API_KEY"))
	tools.NewLichessTools(s, api).RegisterTools()

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
