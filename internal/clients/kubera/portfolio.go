package kubera

import (
	"context"
	"encoding/json"
	"fmt"
)

// Portfolio represents a Kubera portfolio containing its ID, name, assets, and debts.
type Portfolio struct {
	Id     string           `json:"id"`
	Name   string           `json:"name"`
	Assets []*AssetPosition `json:"asset"`
	Debts  []*DebtPosition  `json:"debt"`
}

// getPortfolioResponse wraps the JSON payload returned by the Kubera GetPortfolio API,
// including the Data field on success or an error code on failure.
type getPortfolioResponse struct {
	Data      *Portfolio `json:"data"`
	ErrorCode int64      `json:"errorCode"`
}

// GetPortfolio retrieves the configured portfolio from the Kubera API endpoint,
// validates the response code, and unmarshals the result into a Portfolio struct.
func (c *client) GetPortfolio(ctx context.Context) (*Portfolio, error) {
	data, err := c.get(ctx, fmt.Sprintf("/v3/data/portfolio/%s", c.portfolioId))
	if err != nil {
		return nil, fmt.Errorf("failed to call Kubera API: %w", err)
	}

	var response getPortfolioResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	if response.ErrorCode > 0 {
		return nil, fmt.Errorf("Kubera API error: %d", response.ErrorCode)
	}

	return response.Data, nil
}
