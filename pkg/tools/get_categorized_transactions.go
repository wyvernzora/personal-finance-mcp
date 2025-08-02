package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	ds "github.com/wyvernzora/personal-finance-mcp/pkg/datasource"
)

type GetCategorizedTransactionsInput struct {
	ds.DateRange
}

func GetCategorizedTransactionsTool(ds ds.GetCategorizedTransactionsFunc) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("get_categorized_transactions",
			mcp.WithDescription("Get full transaction list for the specified date range, organized by categories. "),
			mcp.WithString("start_date",
				mcp.Description("Inclusive start date of the interval to list transactions for, formatted like YYYY-MM-DD"),
				mcp.Pattern("[0-9]{4}-[0-9]{2}-[0-9]{2}"),
				mcp.Required(),
			),
			mcp.WithString("end_date",
				mcp.Description("Inclusive end date of the interval to list transactions for, formatted like YYYY-MM-DD"),
				mcp.Pattern("[0-9]{4}-[0-9]{2}-[0-9]{2}"),
				mcp.Required(),
			),
		),
		Handler: mcp.NewTypedToolHandler(
			func(ctx context.Context, _ mcp.CallToolRequest, input GetCategorizedTransactionsInput) (*mcp.CallToolResult, error) {
				result, err := ds(ctx, input.DateRange)
				if err != nil {
					return mcp.NewToolResultErrorFromErr("datasource error", err), err
				}
				return mcp.NewToolResultStructuredOnly(result), nil
			},
		),
	}
}
