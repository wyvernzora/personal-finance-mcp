package types

// Portfolio holds a snapshot of financial positions for a user.
// It tracks net worth, aggregated asset and debt totals, and individual asset/debt entries.
type Portfolio struct {
	NetWorth    Money            `json:"net_worth"`
	TotalAssets Money            `json:"total_assets"`
	TotalDebts  Money            `json:"total_debts"`
	Assets      []*AssetPosition `json:"assets"`
	Debts       []*DebtPosition  `json:"debts"`
}

// NewPortfolio initializes and returns an empty Portfolio with zero balances and no positions.
func NewPortfolio() *Portfolio {
	return &Portfolio{
		NetWorth:    0,
		TotalAssets: 0,
		TotalDebts:  0,
		Assets:      make([]*AssetPosition, 0),
		Debts:       make([]*DebtPosition, 0),
	}
}

// AddAsset adds an AssetPosition to the portfolio, and updates TotalAssets and NetWorth accordingly.
func (p *Portfolio) AddAsset(asset *AssetPosition) {
	p.Assets = append(p.Assets, asset)
	p.TotalAssets += asset.Value
	p.NetWorth += asset.Value
}

// AddDebt adds a DebtPosition to the portfolio, and updates TotalDebts and NetWorth accordingly.
func (p *Portfolio) AddDebt(debt *DebtPosition) {
	p.Debts = append(p.Debts, debt)
	p.TotalDebts += debt.Value
	p.NetWorth -= debt.Value
}
