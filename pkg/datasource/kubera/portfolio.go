package kubera

import (
	"context"
	"iter"
	"log"
	"maps"

	"github.com/bobg/seqs"
	"github.com/wyvernzora/personal-finance-mcp/internal/clients/kubera"
	ds "github.com/wyvernzora/personal-finance-mcp/pkg/datasource"
	"github.com/wyvernzora/personal-finance-mcp/pkg/types"
)

// GetPortfolio fetches portfolio data using the Kubera client stored in context.
// It retrieves raw Kubera assets and debts, transforms them into domain AssetPosition and DebtPosition types,
// and aggregates them into a types.Portfolio.
var GetPortfolio ds.GetPortfolioFunc = func(ctx context.Context) (*types.Portfolio, error) {
	client := kubera.FromContext(ctx)

	// Fetch raw data
	kbPortfolio, err := client.GetPortfolio(ctx)
	if err != nil {
		return nil, err
	}

	// Start processing data
	portfolio := types.NewPortfolio()
	log.Printf("%v", portfolio)
	for asset := range constructAssets(kbPortfolio.Assets) {
		portfolio.AddAsset(asset)
	}
	for debt := range constructDebts(kbPortfolio.Debts) {
		portfolio.AddDebt(debt)
	}

	return portfolio, nil
}

// constructAssets converts Kubera asset positions into a sequence of domain AssetPosition.
// It tombstones non-leaf nodes, filters out parent positions, and annotates each asset with metadata.
func constructAssets(kbAssets []*kubera.AssetPosition) iter.Seq[*types.AssetPosition] {
	assets := make(map[string]*types.AssetPosition)
	for _, kbAsset := range kbAssets {
		// We only care about leaf nodes, so remove or tombstone the parent
		if kbAsset.Parent != nil {
			assets[kbAsset.Parent.Id] = nil
		}

		// Skip if this asset is already tombstoned
		if val, ok := assets[kbAsset.Id]; ok && val == nil {
			continue
		}

		// Construct asset and put into the map
		asset := types.NewAssetPosition(
			kbAsset.Name,
			kbAsset.Ticker,
			determineAssetType(kbAsset),
			kbAsset.AssetClass,
			kbAsset.Value.Amount)
		asset.Description = kbAsset.Description

		asset.Annotate("liquidity", kbAsset.Liquidity)
		asset.Annotate("asset_class", kbAsset.AssetClass)
		asset.Annotate("investable", kbAsset.Investable)
		if kbAsset.Note != "" {
			asset.Annotate("note", kbAsset.Note)
		}

		assets[kbAsset.Id] = asset
	}
	return seqs.Filter(maps.Values(assets), func(v *types.AssetPosition) bool { return v != nil })
}

// determineAssetType maps Kubera asset types and subtypes to standardized domain asset type names.
func determineAssetType(kbAsset *kubera.AssetPosition) string {
	switch {
	case kbAsset.Type == "bank":
		return "cash"
	case kbAsset.Type == "investment":
		return kbAsset.Subtype
	case kbAsset.Type == "other" && kbAsset.Subtype == "home":
		return "real estate"
	default:
		return "unknown"
	}
}

// constructDebts converts Kubera debt positions into a sequence of domain DebtPosition.
// It tombstones non-leaf nodes, filters out parent positions, and returns leaf debt entries.
func constructDebts(kbDebts []*kubera.DebtPosition) iter.Seq[*types.DebtPosition] {
	debts := make(map[string]*types.DebtPosition)
	for _, kbDebt := range kbDebts {
		// We only care about leaf nodes, so remove or tombstone the parent
		if kbDebt.Parent != nil {
			debts[kbDebt.Parent.Id] = nil
		}

		// Skip if this asset is already tombstoned
		if val, ok := debts[kbDebt.Id]; ok && val == nil {
			continue
		}

		// Construct debt and put into the map
		debt := types.NewDebtPosition(
			kbDebt.Name,
			kbDebt.Type,
			kbDebt.Value.Amount)
		debt.Description = kbDebt.Description
		if kbDebt.Note != "" {
			debt.Annotate("note", kbDebt.Note)
		}

		debts[kbDebt.Id] = debt
	}
	return seqs.Filter(maps.Values(debts), func(v *types.DebtPosition) bool { return v != nil })
}
