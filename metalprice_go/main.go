package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	metalprice "github.com/petechu/metalprice-mcp/internal"
)

func main() {
	s := server.NewMCPServer(
		"Metal Price MCP",
		"0.0.1",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	apiKey := os.Getenv("METALPRICE_API_KEY")
	if apiKey == "" {
		log.Fatal("METALPRICE_API_KEY environment variable is not set")
	}

	mcp := metalprice.NewMetalPriceMcp(s, apiKey)

	if err := server.ServeStdio(mcp.Server); err != nil {
		log.Fatalf("Error starting MCP server: %v\n", err)
	}
}
