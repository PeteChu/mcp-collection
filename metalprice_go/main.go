package main

import (
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/server"
	metalprice "github.com/petechu/metalprice-mcp/internal"
)

func main() {
	fmt.Println("Starting MCP server...")

	s := server.NewMCPServer(
		"Metal Price MCP",
		"0.0.1",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	mcp := metalprice.NewMetalPriceMcp(s, "")

	if err := server.ServeStdio(mcp.Server); err != nil {
		log.Fatalf("Error starting MCP server: %v\n", err)
	}
}
