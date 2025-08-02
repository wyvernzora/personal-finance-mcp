package lunchmoney

import (
	"context"
	"net/http"
	"testing"

	lmapi "github.com/wyvernzora/personal-finance-mcp/internal/clients/lunch_money"
	ds "github.com/wyvernzora/personal-finance-mcp/pkg/datasource"
	"github.com/wyvernzora/personal-finance-mcp/pkg/types"
)

// fakeClient implements the LunchMoney API client interface for testing.
type fakeClient struct {
	cats lmapi.Categories
	tags lmapi.Tags
	txs  lmapi.Transactions
}

func (f *fakeClient) ListCategories(ctx context.Context) (lmapi.Categories, error) {
	return f.cats, nil
}
func (f *fakeClient) ListTags(ctx context.Context) (lmapi.Tags, error) {
	return f.tags, nil
}
func (f *fakeClient) ListTransactions(ctx context.Context, startDate, endDate string) (lmapi.Transactions, error) {
	return f.txs, nil
}

// contextWithClient returns a context with the fake client injected.
func contextWithClient(c lmapi.Client) context.Context {
	ctx := context.Background()
	inject := lmapi.WithLunchMoneyClient(c)
	return inject(ctx, &http.Request{})
}

func TestBuildTransaction_Success(t *testing.T) {
	raw := &lmapi.Transaction{
		Date:   "2021-02-03",
		Payee:  "TestPayee",
		Amount: 1234,
		Notes:  "Some notes",
	}
	tx, err := buildTransaction(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got := tx.Payee; got != "TestPayee" {
		t.Errorf("unexpected Payee: %q", got)
	}
	if got := tx.Amount; got != 1234 {
		t.Errorf("unexpected Amount: %d", got)
	}
	if got := tx.Description; got != "Some notes" {
		t.Errorf("unexpected Description: %q", got)
	}
	if got := tx.Date.String(); got != "2021-02-03" {
		t.Errorf("unexpected Date: %q", got)
	}
}

func TestBuildTransaction_InvalidDate(t *testing.T) {
	raw := &lmapi.Transaction{Date: "not-a-date", Payee: "X", Amount: 0}
	if _, err := buildTransaction(raw); err == nil {
		t.Errorf("expected error for invalid date, got nil")
	}
}

func TestAddTransactionTag_SkipsArchived(t *testing.T) {
	tx := types.NewTransaction(types.Date{}, "p", 0)
	archived := &lmapi.Tag{Id: 1, Name: "old", Description: "desc", IsArchived: true}
	addTransactionTag(tx, archived)
	if len(tx.Annotations) != 0 {
		t.Errorf("expected no annotations for archived tag, got %v", tx.Annotations)
	}
}

func TestAddTransactionTag_AnnotationKeyValue(t *testing.T) {
	tx := types.NewTransaction(types.Date{}, "p", 0)
	tag := &lmapi.Tag{Id: 2, Name: "food", Description: "lunch", IsArchived: false}
	addTransactionTag(tx, tag)
	wantKey := "tag:2"
	wantVal := "food: lunch"
	if got, ok := tx.Annotations[wantKey]; !ok || got != wantVal {
		t.Errorf("annotation = %v, want %q", tx.Annotations, wantVal)
	}
}

func TestGetOrCreateCategoryByName_Idempotent(t *testing.T) {
	parent := types.NewCategory("root")
	c1 := getOrCreateCategoryByName(parent, "A")
	c2 := getOrCreateCategoryByName(parent, "A")
	if c1 != c2 {
		t.Errorf("expected same category instance, got different")
	}
	if len(parent.Subcategories) != 1 {
		t.Errorf("expected 1 subcategory, got %d", len(parent.Subcategories))
	}
}

func TestGetOrCreateCategory_CopiesDescription(t *testing.T) {
	root := types.NewCategory("root")
	lmcat := &lmapi.Category{Id: 3, Name: "X", Description: "hello"}
	c := getOrCreateCategory(root, lmcat)
	if c.Description != "hello" {
		t.Errorf("Description = %q, want %q", c.Description, "hello")
	}
}

func TestGetCategorizedTransactions_NoCategory_ExpenseUncategorized(t *testing.T) {
	raw := &lmapi.Transaction{Date: "2022-01-01", Payee: "E", Amount: 50, CategoryId: 0}
	client := &fakeClient{cats: lmapi.Categories{}, tags: lmapi.Tags{}, txs: lmapi.Transactions{raw}}
	ctx := contextWithClient(client)
	start, _ := types.ParseDate("2022-01-01")
	end, _ := types.ParseDate("2022-01-02")
	res, err := GetCategorizedTransactions(ctx, ds.DateRange{StartDate: start, EndDate: end})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	exp := res.Expenses
	if len(exp.Subcategories) != 1 || exp.Subcategories[0].Name != "Uncategorized" {
		t.Errorf("Expenses.Subcategories = %v, want one 'Uncategorized'", exp.Subcategories)
	}
	unc := exp.Subcategories[0]
	if unc.TotalAmount != 50 {
		t.Errorf("Uncategorized.TotalAmount = %v, want 50", unc.TotalAmount)
	}
}

func TestGetCategorizedTransactions_IncomeWithCategory(t *testing.T) {
	raw := &lmapi.Transaction{Date: "2022-02-01", Payee: "I", Amount: 200, IsIncome: true, CategoryId: 10}
	cat := &lmapi.Category{Id: 10, Name: "Salary", Description: ""}
	client := &fakeClient{cats: lmapi.Categories{10: cat}, tags: lmapi.Tags{}, txs: lmapi.Transactions{raw}}
	ctx := contextWithClient(client)
	start, _ := types.ParseDate("2022-02-01")
	end, _ := types.ParseDate("2022-02-02")
	res, err := GetCategorizedTransactions(ctx, ds.DateRange{StartDate: start, EndDate: end})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inc := res.Income
	if len(inc.Subcategories) != 1 || inc.Subcategories[0].Name != "Salary" {
		t.Errorf("Income.Subcategories = %v, want one 'Salary'", inc.Subcategories)
	}
}

func TestGetCategorizedTransactions_MissingCategoryGroup_ErrorAnnotation(t *testing.T) {
	raw := &lmapi.Transaction{Date: "2022-03-01", Payee: "G", Amount: 75, CategoryGroupId: 99, CategoryId: 3}
	cat := &lmapi.Category{Id: 3, Name: "Foo", Description: ""}
	client := &fakeClient{cats: lmapi.Categories{3: cat}, tags: lmapi.Tags{}, txs: lmapi.Transactions{raw}}
	ctx := contextWithClient(client)
	start, _ := types.ParseDate("2022-03-01")
	end, _ := types.ParseDate("2022-03-02")
	res, _ := GetCategorizedTransactions(ctx, ds.DateRange{StartDate: start, EndDate: end})
	exp := res.Expenses.Subcategories[0]
	if exp.Name != "Uncategorized" {
		t.Errorf("got %q, want 'Uncategorized'", exp.Name)
	}
	if val, ok := exp.Transactions[0].Annotations["category_error"]; !ok || val == "" {
		t.Errorf("expected category_error annotation, got %v", exp.Transactions[0].Annotations)
	}
}

func TestGetCategorizedTransactions_MissingCategoryId_ErrorAnnotation(t *testing.T) {
	raw := &lmapi.Transaction{Date: "2022-06-01", Payee: "Y", Amount: 123, CategoryId: 42}
	client := &fakeClient{cats: lmapi.Categories{}, tags: lmapi.Tags{}, txs: lmapi.Transactions{raw}}
	ctx := contextWithClient(client)
	start, _ := types.ParseDate("2022-06-01")
	end, _ := types.ParseDate("2022-06-02")
	res, _ := GetCategorizedTransactions(ctx, ds.DateRange{StartDate: start, EndDate: end})
	unc := res.Expenses.Subcategories[0]
	if val, ok := unc.Transactions[0].Annotations["category_error"]; !ok || val == "" {
		t.Errorf("expected category_error annotation, got %v", unc.Transactions[0].Annotations)
	}
}

func TestGetCategorizedTransactions_CategoryGroupMatchingName(t *testing.T) {
	raw := &lmapi.Transaction{Date: "2022-07-01", Payee: "Z", Amount: 150, CategoryGroupId: 1, CategoryId: 2}
	group := &lmapi.Category{Id: 1, Name: "Expenses", Description: ""}
	cat := &lmapi.Category{Id: 2, Name: "Food", Description: "desc"}
	client := &fakeClient{cats: lmapi.Categories{1: group, 2: cat}, tags: lmapi.Tags{}, txs: lmapi.Transactions{raw}}
	ctx := contextWithClient(client)
	start, _ := types.ParseDate("2022-07-01")
	end, _ := types.ParseDate("2022-07-02")
	res, _ := GetCategorizedTransactions(ctx, ds.DateRange{StartDate: start, EndDate: end})
	sub := res.Expenses.Subcategories
	if len(sub) != 1 || sub[0].Name != "Food" || sub[0].Description != "desc" {
		t.Errorf("got %#+v, want one Food with desc", sub)
	}
}

func TestGetCategorizedTransactions_ExcludeFromBudget(t *testing.T) {
	raw := &lmapi.Transaction{Date: "2022-05-01", Payee: "X", Amount: 99, ExcludeFromBudget: true, CategoryId: 0}
	client := &fakeClient{cats: lmapi.Categories{}, tags: lmapi.Tags{}, txs: lmapi.Transactions{raw}}
	ctx := contextWithClient(client)
	start, _ := types.ParseDate("2022-05-01")
	end, _ := types.ParseDate("2022-05-02")
	res, _ := GetCategorizedTransactions(ctx, ds.DateRange{StartDate: start, EndDate: end})
	ign := res.Ignored.Subcategories
	if len(ign) != 1 || ign[0].Name != "Uncategorized" {
		t.Errorf("Ignored = %v, want one Uncategorized", ign)
	}
}

func TestGetCategorizedTransactions_ExcludeFromTotals(t *testing.T) {
	raw := &lmapi.Transaction{Date: "2022-08-01", Payee: "O", Amount: 42, ExcludeFromTotals: true, CategoryId: 0}
	client := &fakeClient{cats: lmapi.Categories{}, tags: lmapi.Tags{}, txs: lmapi.Transactions{raw}}
	ctx := contextWithClient(client)
	start, _ := types.ParseDate("2022-08-01")
	end, _ := types.ParseDate("2022-08-02")
	res, _ := GetCategorizedTransactions(ctx, ds.DateRange{StartDate: start, EndDate: end})
	ign := res.Ignored.Subcategories
	if len(ign) != 1 || ign[0].TotalAmount != 42 {
		t.Errorf("ExcludeFromTotals not applied, got %v", ign)
	}
}

func TestGetCategorizedTransactions_AttachesTags(t *testing.T) {
	raw := &lmapi.Transaction{
		Date:       "2022-04-01",
		Payee:      "T",
		Amount:     30,
		CategoryId: 0,
		Tags: []*struct {
			Id int64 `json:"id"`
		}{{Id: 5}, {Id: 6}},
	}
	client := &fakeClient{
		cats: lmapi.Categories{},
		tags: lmapi.Tags{
			5: {Id: 5, Name: "tag5", Description: "d5", IsArchived: false},
			6: {Id: 6, Name: "tag6", Description: "d6", IsArchived: true},
		},
		txs: lmapi.Transactions{raw},
	}
	ctx := contextWithClient(client)
	start, _ := types.ParseDate("2022-04-01")
	end, _ := types.ParseDate("2022-04-02")
	res, _ := GetCategorizedTransactions(ctx, ds.DateRange{StartDate: start, EndDate: end})
	unc := res.Expenses.Subcategories[0]
	if _, has := unc.Transactions[0].Annotations["tag:6"]; has {
		t.Errorf("archived tag annotated: %v", unc.Transactions[0].Annotations)
	}
	if val, has := unc.Transactions[0].Annotations["tag:5"]; !has || val != "tag5: d5" {
		t.Errorf("tag5 missing: %v", unc.Transactions[0].Annotations)
	}
}

func TestGetCategorizedTransactions_MissingTagSkipped(t *testing.T) {
	raw := &lmapi.Transaction{
		Date:       "2022-09-01",
		Payee:      "M",
		Amount:     5,
		CategoryId: 0,
		Tags: []*struct {
			Id int64 `json:"id"`
		}{{Id: 123}},
	}
	client := &fakeClient{cats: lmapi.Categories{}, tags: lmapi.Tags{}, txs: lmapi.Transactions{raw}}
	ctx := contextWithClient(client)
	start, _ := types.ParseDate("2022-09-01")
	end, _ := types.ParseDate("2022-09-02")
	res, _ := GetCategorizedTransactions(ctx, ds.DateRange{StartDate: start, EndDate: end})
	if len(res.Expenses.Subcategories[0].Transactions[0].Annotations) != 0 {
		t.Errorf("expected no annotations for missing tag, got %v", res.Expenses.Subcategories[0].Transactions[0].Annotations)
	}
}
