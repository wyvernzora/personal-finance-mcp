package datasource

import (
	"context"

	"github.com/wyvernzora/personal-finance-mcp/pkg/types"
)

// DateRange defines the inclusive start and end dates for querying data.
type DateRange struct {
	StartDate types.Date `json:"start_date"`
	EndDate   types.Date `json:"end_date"`
}

// GetCategorizedTransactionsFunc is the signature of a data source method that fetches, categorizes,
// and returns transactions within the specified DateRange.
type GetCategorizedTransactionsFunc func(ctx context.Context, interval DateRange) (*types.Categories, error)

// GetPortfolioFunc is the signature of a data source method that retrieves the current portfolio,
// including all asset and debt positions.
type GetPortfolioFunc func(ctx context.Context) (*types.Portfolio, error)
