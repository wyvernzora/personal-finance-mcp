package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	bindAddr := os.Getenv("BIND_ADDRESS")
	if bindAddr == "" {
		bindAddr = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	httpServer := server.NewStreamableHTTPServer(createMCPServer())
	log.Printf("HTTP server listening on %s:%s/mcp", bindAddr, port)
	if err := httpServer.Start(bindAddr + ":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func createMCPServer() *server.MCPServer {
	mcpServer := server.NewMCPServer(
		"personal-finance-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithInstructions(
			"This server provides tools to retrieve information about a user's personal finances, such as "+
				"spending transactions, asset holdings, net worth etc",
		),
	)
	return mcpServer
}
