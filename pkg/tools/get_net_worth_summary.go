package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	ds "github.com/wyvernzora/personal-finance-mcp/pkg/datasource"
)

func GetNetWorthSummary(ds ds.GetPortfolioFunc) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("get_net_worth_summary",
			mcp.WithDescription("Get a summary of user's net worth, including all asset holdings, debts and their respective values."),
		),
		Handler: func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := ds(ctx)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("datasource error", err), err
			}
			return mcp.NewToolResultStructuredOnly(result), nil
		},
	}
}
