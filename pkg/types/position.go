package types

// Position represents a financial position (asset or debt) with a name,
// type, and monetary value. It also includes annotations for system metadata.
type Position struct {
	AnnotatedObject
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value Money  `json:"value"`
}

// AssetPosition represents an asset holding, extending Position with a ticker symbol.
type AssetPosition struct {
	Position
	Ticker string `json:"ticker"`
}

// NewAssetPosition constructs a new AssetPosition with the provided name, ticker,
// asset type, asset class, and value, and initializes its annotations map.
func NewAssetPosition(name, ticker, assetType, assetClass string, value Money) *AssetPosition {
	return &AssetPosition{
		Position: Position{
			Name:            name,
			Type:            assetType,
			Value:           value,
			AnnotatedObject: NewAnnotatedObject(),
		},
		Ticker: ticker,
	}
}

// DebtPosition represents a liability or debt position, extending Position.
type DebtPosition struct {
	Position
}

// NewDebtPosition constructs a new DebtPosition with the provided name, type,
// and value, and initializes its annotations map.
func NewDebtPosition(name, debtType string, value Money) *DebtPosition {
	return &DebtPosition{
		Position: Position{
			Name:            name,
			Type:            debtType,
			Value:           value,
			AnnotatedObject: NewAnnotatedObject(),
		},
	}
}
