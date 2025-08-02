package kubera

import (
	"context"
	"testing"

	clients "github.com/wyvernzora/personal-finance-mcp/internal/clients/kubera"
	"github.com/wyvernzora/personal-finance-mcp/pkg/types"
)

// TestDetermineAssetType verifies mapping of Kubera asset types and subtypes to domain types.
func TestDetermineAssetType(t *testing.T) {
	cases := []struct {
		typ, subtype, want string
	}{
		{"bank", "", "cash"},
		{"investment", "stock", "stock"},
		{"other", "home", "real estate"},
		{"other", "car", "unknown"},
		{"foo", "bar", "unknown"},
	}

	for _, c := range cases {
		kb := &clients.AssetPosition{
			Position: clients.Position{
				Type:    c.typ,
				Subtype: c.subtype,
			},
		}
		got := determineAssetType(kb)
		if got != c.want {
			t.Errorf("determineAssetType(%q, %q) = %q; want %q", c.typ, c.subtype, got, c.want)
		}
	}
}

// TestConstructAssets_Tombstone ensures parent assets are tombstoned when children exist.
func TestConstructAssets_Tombstone(t *testing.T) {
	parent := &clients.AssetPosition{
		Position: clients.Position{
			Id:     "p1",
			Name:   "ParentAsset",
			Ticker: "PARENT",
			Value:  clients.Value{Amount: 1000, Currency: "USD"},
		},
		Investable: "yes",
		Liquidity:  "high",
		AssetClass: "equity",
	}
	child := &clients.AssetPosition{
		Position: clients.Position{
			Id:      "c1",
			Name:    "ChildAsset",
			Type:    "investment",
			Subtype: "stock",
			Ticker:  "CHILD",
			Value:   clients.Value{Amount: 2500, Currency: "USD"},
			Parent:  &clients.ParentPosition{Id: "p1", Name: "ParentAsset"},
		},
		Investable: "no",
		Liquidity:  "low",
		AssetClass: "bond",
	}

	seq := constructAssets([]*clients.AssetPosition{parent, child})
	var out []*types.AssetPosition
	for a := range seq {
		out = append(out, a)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 leaf asset, got %d", len(out))
	}
	got := out[0]

	if got.Name != "ChildAsset" {
		t.Errorf("Name = %q; want %q", got.Name, "ChildAsset")
	}
	if got.Ticker != "CHILD" {
		t.Errorf("Ticker = %q; want %q", got.Ticker, "CHILD")
	}
	if got.Type != "stock" {
		t.Errorf("Type = %q; want %q", got.Type, "stock")
	}
	if got.Value != types.Money(2500) {
		t.Errorf("Value = %v; want %v", got.Value, types.Money(2500))
	}

	wantAnn := map[string]string{
		"liquidity":   "low",
		"asset_class": "bond",
		"investable":  "no",
	}
	for k, v := range wantAnn {
		if got.Annotations[k] != v {
			t.Errorf("Annotations[%q] = %q; want %q", k, got.Annotations[k], v)
		}
	}
}

// TestConstructAssets_NoTombstone ensures assets without parent-child relationships are all returned.
func TestConstructAssets_NoTombstone(t *testing.T) {
	a1 := &clients.AssetPosition{
		Position: clients.Position{Id: "a1", Name: "A1", Ticker: "A1", Value: clients.Value{Amount: 100, Currency: "USD"}},
	}
	a2 := &clients.AssetPosition{
		Position: clients.Position{Id: "a2", Name: "A2", Ticker: "A2", Value: clients.Value{Amount: 200, Currency: "USD"}},
	}

	seq := constructAssets([]*clients.AssetPosition{a1, a2})
	var names []string
	for a := range seq {
		names = append(names, a.Name)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(names))
	}
	set := map[string]bool{names[0]: true, names[1]: true}
	if !set["A1"] || !set["A2"] {
		t.Errorf("unexpected assets %v; want A1 and A2", names)
	}
}

// TestConstructDebts_Tombstone ensures parent debts are tombstoned when children exist.
func TestConstructDebts_Tombstone(t *testing.T) {
	parent := &clients.DebtPosition{
		Position: clients.Position{
			Id:    "pD",
			Name:  "ParentDebt",
			Type:  "loan",
			Value: clients.Value{Amount: 500, Currency: "USD"},
		},
	}
	child := &clients.DebtPosition{
		Position: clients.Position{
			Id:     "cD",
			Name:   "ChildDebt",
			Type:   "mortgage",
			Value:  clients.Value{Amount: 1500, Currency: "USD"},
			Parent: &clients.ParentPosition{Id: "pD", Name: "ParentDebt"},
		},
	}

	seq := constructDebts([]*clients.DebtPosition{parent, child})
	var out []*types.DebtPosition
	for d := range seq {
		out = append(out, d)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 leaf debt, got %d", len(out))
	}
	got := out[0]
	if got.Name != "ChildDebt" {
		t.Errorf("Name = %q; want %q", got.Name, "ChildDebt")
	}
	if got.Type != "mortgage" {
		t.Errorf("Type = %q; want %q", got.Type, "mortgage")
	}
	if got.Value != types.Money(1500) {
		t.Errorf("Value = %v; want %v", got.Value, types.Money(1500))
	}
}

// TestConstructDebts_NoTombstone ensures debts without parent-child relationships are all returned.
func TestConstructDebts_NoTombstone(t *testing.T) {
	d1 := &clients.DebtPosition{
		Position: clients.Position{Id: "d1", Name: "D1", Type: "cc", Value: clients.Value{Amount: 300, Currency: "USD"}},
	}
	d2 := &clients.DebtPosition{
		Position: clients.Position{Id: "d2", Name: "D2", Type: "loan", Value: clients.Value{Amount: 700, Currency: "USD"}},
	}

	seq := constructDebts([]*clients.DebtPosition{d1, d2})
	var names []string
	for d := range seq {
		names = append(names, d.Name)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 debts, got %d", len(names))
	}
	set := map[string]bool{names[0]: true, names[1]: true}
	if !set["D1"] || !set["D2"] {
		t.Errorf("unexpected debts %v; want D1 and D2", names)
	}
}

// TestGetPortfolio_Success verifies GetPortfolio properly transforms assets and debts.
func TestGetPortfolio_Success(t *testing.T) {
	// Inject fake client into context
	ctx := clients.WithKuberaClient(&fakeClient{})(context.Background(), nil)
	// Call under test
	p, err := GetPortfolio(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Validate totals
	if p.TotalAssets != types.Money(100) {
		t.Errorf("TotalAssets = %v; want %v", p.TotalAssets, types.Money(100))
	}
	if p.TotalDebts != types.Money(40) {
		t.Errorf("TotalDebts = %v; want %v", p.TotalDebts, types.Money(40))
	}
	if p.NetWorth != types.Money(60) {
		t.Errorf("NetWorth = %v; want %v", p.NetWorth, types.Money(60))
	}
	// Validate contents
	if len(p.Assets) != 1 || p.Assets[0].Name != "Asset1" {
		t.Errorf("Assets = %v; want single Asset1", p.Assets)
	}
	if len(p.Debts) != 1 || p.Debts[0].Name != "Debt1" {
		t.Errorf("Debts = %v; want single Debt1", p.Debts)
	}
}

// TestGetPortfolio_Panic_NoClient verifies that GetPortfolio panics without a client in context.
func TestGetPortfolio_Panic_NoClient(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when no Kubera client in context")
		}
	}()
	_, _ = GetPortfolio(context.Background())
}
