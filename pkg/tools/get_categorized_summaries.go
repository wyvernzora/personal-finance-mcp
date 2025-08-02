package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	ds "github.com/wyvernzora/personal-finance-mcp/pkg/datasource"
	"github.com/wyvernzora/personal-finance-mcp/pkg/types"
)

type GetCategorizedSummariesInput struct {
	ds.DateRange
}

func GetCategorizedSummariesTool(ds ds.GetCategorizedTransactionsFunc) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("get_categorized_summaries",
			mcp.WithDescription(
				"Get spending summary by category for the specified date range. Does NOT include actual transaction list. "+
					"Use when assessing long term trends where drilling into individual transactions is not necessary. "+
					"get_categorized_transactions tool can provide full list of transactions if needed",
			),
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
			func(ctx context.Context, _ mcp.CallToolRequest, input GetCategorizedSummariesInput) (*mcp.CallToolResult, error) {
				result, err := ds(ctx, input.DateRange)
				if err != nil {
					return mcp.NewToolResultErrorFromErr("datasource error", err), err
				}
				removeTransactions(result.Income)
				removeTransactions(result.Expenses)
				removeTransactions(result.Ignored)
				return mcp.NewToolResultStructuredOnly(result), nil
			},
		),
	}
}

func removeTransactions(cat *types.Category) {
	cat.Transactions = nil
	for _, sub := range cat.Subcategories {
		removeTransactions(sub)
	}
}
