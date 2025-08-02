package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wyvernzora/personal-finance-mcp/pkg/types"
)

// Transactions is a collection of Transaction pointers returned by ListTransactions.
type Transactions []*Transaction

// Transaction represents a single financial entry returned by the Lunch Money API.
// It includes metadata such as date, payee, category/group IDs and names, amounts, notes, and tags.
type Transaction struct {
	Id                   int64       `json:"id"`
	Date                 string      `json:"date"`
	Amount               types.Money `json:"to_base"`
	Payee                string      `json:"payee"`
	OriginalPayee        string      `json:"original_name"`
	CategoryId           int64       `json:"category_id"`
	CategoryName         string      `json:"category_name"`
	CategoryGroupId      int64       `json:"category_group_id"`
	CategoryGroupName    string      `json:"category_group_name"`
	IsIncome             bool        `json:"is_income"`
	ExcludeFromBudget    bool        `json:"exclude_from_budget"`
	ExcludeFromTotals    bool        `json:"exclude_from_totals"`
	Notes                string      `json:"display_notes"`
	RecurringCadence     string      `json:"recurring_cadence,omitempty"`
	RecurringDescription string      `json:"recurring_description,omitempty"`
	Tags                 []*struct {
		Id int64 `json:"id"`
	} `json:"tags"`
}

// listTransactionsResponse wraps the JSON payload from the Lunch Money ListTransactions endpoint,
// containing the slice of transactions and a flag indicating if more pages exist.
type listTransactionsResponse struct {
	Transactions []*Transaction `json:"transactions"`
	HasMore      bool           `json:"has_more"`
}

// ListTransactions retrieves all transactions between startDate and endDate (inclusive) from Lunch Money.
// Dates must be in "YYYY-MM-DD" format. It returns a Transactions slice or an error if the API call or unmarshaling fails.
func (c *client) ListTransactions(ctx context.Context, startDate, endDate string) (Transactions, error) {
	data, err := c.get(ctx, "/v1/transactions", map[string]string{
		"start_date": startDate,
		"end_date":   endDate,
		"limit":      "10000",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call Lunch Money API: %w", err)
	}

	var response listTransactionsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}
	if response.HasMore {
		// 10K ceiling should be sufficient for most ML use cases
		// Implement proper pagination if ever needed
		return nil, fmt.Errorf("too many transactions, try smaller time interval")
	}

	return response.Transactions, nil
}
