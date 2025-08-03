package kubera

import "github.com/wyvernzora/personal-finance-mcp/pkg/types"

// Value wraps an amount and currency for a Kubera position.
type Value struct {
	Amount   types.Money `json:"amount"`
	Currency string      `json:"currency"`
}

// Geography holds the country and region metadata for a Kubera position.
type Geography struct {
	Country string `json:"country"`
	Region  string `json:"region"`
}

// ParentPosition identifies the parent of a position by its ID and name.
type ParentPosition struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Position represents a single portfolio entry from Kubera, including its ID, name, value, and hierarchy.
type Position struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Note        string          `json:"note"`
	Value       Value           `json:"value"`
	Ticker      string          `json:"ticker"`
	Type        string          `json:"type"`
	Subtype     string          `json:"subType"`
	Parent      *ParentPosition `json:"parent"`
}

// AssetPosition extends Position with asset-specific fields such as investability, liquidity, and asset class.
type AssetPosition struct {
	Position
	Investable string     `json:"investable"`
	Liquidity  string     `json:"liquidity"`
	AssetClass string     `json:"assetClass"`
	Geography  *Geography `json:"geography"`
}

// DebtPosition extends Position to represent liabilities or debts in a Kubera portfolio.
type DebtPosition struct {
	Position
}
