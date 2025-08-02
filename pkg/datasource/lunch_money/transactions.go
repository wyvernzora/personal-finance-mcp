package lunchmoney

import (
	"context"
	"fmt"
	"log"

	lmapi "github.com/wyvernzora/personal-finance-mcp/internal/clients/lunch_money"
	ds "github.com/wyvernzora/personal-finance-mcp/pkg/datasource"
	"github.com/wyvernzora/personal-finance-mcp/pkg/types"
)

// GetCategorizedTransactions is a DataSource function that fetches transactions from LunchMoney API,
// categorizes them into Income, Expenses, and Ignored buckets, and converts them into domain-specific
// types enriched with annotations from tags and error metadata. It handles missing categories and groups,
// placing uncategorized transactions accordingly and adds error annotations when inconsistencies occur.
var GetCategorizedTransactions ds.GetCategorizedTransactionsFunc = func(ctx context.Context, interval ds.DateRange) (*types.Categories, error) {
	client := lmapi.FromContext(ctx)

	// Grab raw data from LunchMoney
	lmCats, err := client.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	lmTags, err := client.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	lmTxs, err := client.ListTransactions(ctx, interval.StartDate.String(), interval.EndDate.String())
	if err != nil {
		return nil, err
	}

	// Maps to keep track of categories as we add stuff to them
	result := types.NewCategories()

	// Start processing transactions
	for _, lmtx := range lmTxs {
		tx, err := buildTransaction(lmtx)
		if err != nil {
			return nil, err
		}

		// Attach tags as annotations
		for _, tag := range lmtx.Tags {
			lmtag, ok := lmTags[tag.Id]
			if !ok {
				log.Printf("missing tag from LunchMoney response: %d", tag.Id)
				continue
			}
			addTransactionTag(tx, lmtag)
		}

		// Determine which "bucket" does the transaction fall under
		var bucket *types.Category
		switch {
		case lmtx.IsIncome:
			bucket = result.Income
		case lmtx.ExcludeFromBudget || lmtx.ExcludeFromTotals:
			bucket = result.Ignored
		default:
			bucket = result.Expenses
		}

		// Utility function to set transaction as uncategorized
		setAsUncategorized := func() error {
			unc := getOrCreateCategoryByName(bucket, "Uncategorized")
			return unc.AddTransaction(tx)
		}

		// Case 1: no category, put under Uncategorized
		if lmtx.CategoryId == 0 {
			if err := setAsUncategorized(); err != nil {
				return nil, err
			}
			continue
		}
		// Case 2: has category group; make category group the bucket
		if lmtx.CategoryGroupId != 0 {
			lmCatGroup, ok := lmCats[lmtx.CategoryGroupId]
			if !ok {
				log.Printf("missing category group from LunchMoney response: %d\n", lmtx.CategoryGroupId)
				tx.Annotate("category_error", "uncategorized due to invalid category group id")
				if err := setAsUncategorized(); err != nil {
					return nil, err
				}
				continue
			}
			if lmCatGroup.Name != bucket.Name {
				catGroup := getOrCreateCategory(bucket, lmCatGroup)
				bucket = catGroup
			}
		}
		// Add transaction to bucket
		lmCat, ok := lmCats[lmtx.CategoryId]
		if !ok {
			log.Printf("missing category from LunchMoney response: %d\n", lmtx.CategoryId)
			tx.Annotate("category_error", "uncategorized due to invalid category id")
			if err := setAsUncategorized(); err != nil {
				return nil, err
			}
			continue
		}
		cat := getOrCreateCategory(bucket, lmCat)
		if err := cat.AddTransaction(tx); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// buildTransaction converts a LunchMoney API transaction into a domain Transaction type,
// setting date, payee, and amount, and copying over the notes as description.
func buildTransaction(lmtx *lmapi.Transaction) (*types.Transaction, error) {
	date, err := types.ParseDate(lmtx.Date)
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(date, lmtx.Payee, lmtx.Amount)
	tx.Description = lmtx.Notes
	return tx, nil
}

// addTransactionTag annotates a Transaction with the given LunchMoney Tag providing
// key-value metadata. Skips archived tags.
func addTransactionTag(tx *types.Transaction, tag *lmapi.Tag) {
	if tag.IsArchived {
		log.Printf("skipping archived tag: %d - %s\n", tag.Id, tag.Name)
		return
	}
	key := fmt.Sprintf("tag:%d", tag.Id)
	value := fmt.Sprintf("%s: %s", tag.Name, tag.Description)
	tx.Annotate(key, value)
}

// getOrCreateCategory finds or creates a subcategory under the parent Category
// matching the given LunchMoney Category. It also copies the description from LunchMoney.
func getOrCreateCategory(parent *types.Category, lmcat *lmapi.Category) *types.Category {
	cat := getOrCreateCategoryByName(parent, lmcat.Name)
	cat.Description = lmcat.Description
	return cat
}

// getOrCreateCategoryByName returns an existing subcategory by name under the parent,
// or creates and attaches a new Category with that name if none exists.
func getOrCreateCategoryByName(parent *types.Category, name string) *types.Category {
	for _, sc := range parent.Subcategories {
		if sc.Name == name {
			return sc
		}
	}
	newCat := types.NewCategory(name)
	_ = parent.AddSubcategory(newCat)
	return newCat
}
