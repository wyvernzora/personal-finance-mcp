package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// helper to create a dummy transaction with EVE-themed payee and amount
func makeTestTransaction(payee string, amount Money) *Transaction {
	return NewTransaction(Date{}, payee, amount)
}

func TestAddSubcategory_Success(t *testing.T) {
	const (
		rootName = "Gallente Federation"
		subName  = "Minmatar Republic"
	)
	root := NewCategory(rootName)
	sub := NewCategory(subName)

	if err := root.AddSubcategory(sub); err != nil {
		t.Fatalf("unexpected error adding subcategory: %v", err)
	}
	if len(root.Subcategories) != 1 {
		t.Errorf("expected 1 subcategory, got %d", len(root.Subcategories))
	}
	if sub.Parent != root {
		t.Errorf("expected sub.Parent to be root, got %v", sub.Parent)
	}
	// no transactions, so total should remain zero
	if root.TotalAmount != 0 {
		t.Errorf("expected root.TotalAmount 0, got %v", root.TotalAmount)
	}
}

func TestAddSubcategory_ErrorAlreadyParent(t *testing.T) {
	p1 := NewCategory("Amarr Empire")
	p2 := NewCategory("Caldari State")
	sub := NewCategory("Jove Observatory")

	if err := p1.AddSubcategory(sub); err != nil {
		t.Fatalf("setup: unexpected error adding to first parent: %v", err)
	}
	err := p2.AddSubcategory(sub)
	if err == nil {
		t.Fatalf("expected error when adding subcategory with existing parent, got nil")
	}
	want := fmt.Sprintf("cannot move subcategory %s to another parent", sub.Name)
	if !contains(err.Error(), want) {
		t.Errorf("error message %q does not contain %q", err.Error(), want)
	}
}

func TestAddTransaction_Success(t *testing.T) {
	cat := NewCategory("Sisters of EVE")
	txn := makeTestTransaction("Rookie Help", Money(25000)) // 2.5000 ISK

	if err := cat.AddTransaction(txn); err != nil {
		t.Fatalf("unexpected error adding transaction: %v", err)
	}
	if len(cat.Transactions) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(cat.Transactions))
	}
	if txn.Category != cat {
		t.Errorf("expected txn.Category to be %v, got %v", cat, txn.Category)
	}
	if cat.TotalAmount != 25000 {
		t.Errorf("expected TotalAmount 25000, got %v", cat.TotalAmount)
	}
}

func TestAddTransaction_ErrorAlreadyAssigned(t *testing.T) {
	c1 := NewCategory("Mordus Legion")
	c2 := NewCategory("CONCORD")
	txn := makeTestTransaction("Security Service", Money(10000))

	if err := c1.AddTransaction(txn); err != nil {
		t.Fatalf("setup: unexpected error adding transaction to first category: %v", err)
	}
	err := c2.AddTransaction(txn)
	if err == nil {
		t.Fatal("expected error when reassigning transaction, got nil")
	}
	want := "cannot move transaction to another parent"
	if !contains(err.Error(), want) {
		t.Errorf("error message %q does not contain %q", err.Error(), want)
	}
}

// TestRecomputeTotals ensures nested subcategories and transactions sum correctly.
func TestRecomputeTotals(t *testing.T) {
	root := NewCategory("Empire Wallet")
	sec := NewCategory("Faction Wallet")
	sub1 := NewCategory("Amarr Tax")
	sub2 := NewCategory("Minmatar Gift")

	// Build tree: root -> sec -> {sub1, sub2}
	if err := root.AddSubcategory(sec); err != nil {
		t.Fatal(err)
	}
	if err := sec.AddSubcategory(sub1); err != nil {
		t.Fatal(err)
	}
	if err := sec.AddSubcategory(sub2); err != nil {
		t.Fatal(err)
	}
	// Add transactions
	_ = sub1.AddTransaction(makeTestTransaction("Tribute", Money(50000)))   // 5.0000
	_ = sub2.AddTransaction(makeTestTransaction("Aid", Money(20000)))       // 2.0000
	_ = root.AddTransaction(makeTestTransaction("Donations", Money(30000))) // 3.0000

	// Now recomputeTotals on root
	total := root.recomputeTotals()
	// root should be 3 + (5+2) = 10.0000 => 100000
	if total != 100000 {
		t.Errorf("expected total 100000, got %v", total)
	}
	if root.TotalAmount != total {
		t.Errorf("root.TotalAmount %v != returned %v", root.TotalAmount, total)
	}
}

// TestCategories_UnmarshalJSON verifies JSON loading, parent links, and transaction pointers.
func TestCategories_UnmarshalJSON(t *testing.T) {
	jsonData := `
{
  "income": {
    "name": "Income",
    "total_amount": 150.0000,
    "subcategories": [
      {
        "name": "Federation Funds",
        "total_amount": 150.0000,
        "subcategories": [
          {
            "name": "GalNet Ads",
            "total_amount": 50.0000,
            "transactions": [
              { "date": "2023-01-01", "payee": "GalNet", "amount": 25.0000 },
              { "date": "2023-01-02", "payee": "GalNet", "amount": 25.0000 }
            ]
          }
        ],
        "transactions": [
          { "date": "2023-01-03", "payee": "SpectreFleet", "amount": 100.0000 }
        ]
      }
    ]
  },
  "expenses": {},
  "ignored": {}
}`

	var cats Categories
	if err := json.Unmarshal([]byte(jsonData), &cats); err != nil {
		t.Fatalf("UnmarshalJSON error: %v", err)
	}

	// Check top-level income count
	if len(cats.Income.Subcategories) != 1 {
		t.Fatalf("expected 1 income category, got %d", len(cats.Income.Subcategories))
	}
	fed := cats.Income.Subcategories[0]
	if fed.Name != "Federation Funds" {
		t.Errorf("unexpected name %q", fed.Name)
	}
	// Check subcategory
	if len(fed.Subcategories) != 1 {
		t.Fatalf("expected 1 subcategory, got %d", len(fed.Subcategories))
	}
	galnet := fed.Subcategories[0]
	if galnet.Parent != fed {
		t.Errorf("expected galnet.Parent == fed, got %v", galnet.Parent)
	}
	// Check transactions pointers and counts
	if len(galnet.Transactions) != 2 {
		t.Errorf("expected 2 galnet transactions, got %d", len(galnet.Transactions))
	}
	for i, txn := range galnet.Transactions {
		if txn.Category != galnet {
			t.Errorf("txn[%d].Category = %v, want %v", i, txn.Category, galnet)
		}
	}
	if len(fed.Transactions) != 1 {
		t.Errorf("expected 1 root transaction, got %d", len(fed.Transactions))
	}
	// Recompute totals and check against JSON-provided amounts
	got := fed.recomputeTotals()
	want := Money(1500000) // 150.0000
	if got != want {
		t.Errorf("recomputeTotals = %v, want %v", got, want)
	}
}

// contains is a helper for substring checks in tests.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
