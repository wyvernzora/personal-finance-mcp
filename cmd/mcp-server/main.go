package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/wyvernzora/personal-finance-mcp/pkg/datasource/kubera"
	lm "github.com/wyvernzora/personal-finance-mcp/pkg/datasource/lunch_money"
	"github.com/wyvernzora/personal-finance-mcp/pkg/tools"
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

	httpServer := server.NewStreamableHTTPServer(
		createMCPServer(),
		withHTTPContextFuncs(
			lm.InjectCredentialsFromEnvironment(),
			kubera.InjectCredentialsFromEnvironment(),
		),
	)
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

	mcpServer.AddTools(
		tools.GetCategorizedTransactionsTool(lm.GetCategorizedTransactions),
		tools.GetCategorizedSummariesTool(lm.GetCategorizedTransactions),
		tools.GetNetWorthSummary(kubera.GetPortfolio),
	)

	return mcpServer
}

func withHTTPContextFuncs(fns ...server.HTTPContextFunc) server.StreamableHTTPOption {
	return server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
		for _, fn := range fns {
			ctx = fn(ctx, r)
		}
		return ctx
	})
}
