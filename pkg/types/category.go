package types

import (
	"encoding/json"
	"fmt"
)

// Category represents a financial category which may contain nested subcategories
// and associated transactions. It tracks its own total amount as the sum of its
// transactions plus the totals of its subcategories.
type Category struct {
	// AnnotatedObject provides common description fields.
	AnnotatedObject
	// Name is the unique identifier for the category.
	Name string `json:"name"`
	// Description is an optional text provided by the end user for context.
	Description string `json:"description,omitempty"`
	// Parent references the parent category in the hierarchy. It is nil for top‐level categories.
	// It is automatically set by Category.AddSubcategory.
	// To restore links after JSON unmarshaling, call rebuildTree on the root Category.
	Parent *Category `json:"-"`
	// TotalAmount is the aggregated sum of all transactions and subcategory totals.
	TotalAmount Money `json:"total_amount"`
	// Subcategories holds child categories nested under this category.
	Subcategories []*Category `json:"subcategories,omitempty"`
	// Transactions are the individual financial entries assigned to this category.
	Transactions []*Transaction `json:"transactions,omitempty"`
}

// NewCategory constructs and returns a pointer to a Category with the provided name.
// The resulting Category has no parent, zero total amount, and empty slices for
// subcategories and transactions.
func NewCategory(name string) *Category {
	return &Category{
		Name:            name,
		AnnotatedObject: NewAnnotatedObject(),
		TotalAmount:     0,
		Subcategories:   make([]*Category, 0),
		Transactions:    make([]*Transaction, 0),
	}
}

// rebuildTree recursively sets Parent pointers on all subcategories and attaches
// this category as the Category field on its transactions.
func (c *Category) rebuildTree() {
	for _, sub := range c.Subcategories {
		sub.Parent = c
		sub.rebuildTree()
	}
	for _, txn := range c.Transactions {
		txn.Category = c
	}
}

// recomputeTotals walks the entire subtree rooted at this Category, recalculates
// each node’s TotalAmount, and returns the computed total for this node.
func (c *Category) recomputeTotals() Money {
	var total Money
	for _, sub := range c.Subcategories {
		total += sub.recomputeTotals()
	}
	for _, txn := range c.Transactions {
		total += txn.Amount
	}
	c.TotalAmount = total
	return total
}

func (c *Category) addToTotalAmount(amt Money) {
	c.TotalAmount += amt
	if c.Parent != nil {
		c.Parent.addToTotalAmount(amt)
	}
}

// AddSubcategory attaches the provided subcategory beneath this Category,
// sets its Parent pointer, and recomputes the total amounts up the tree.
// It returns an error if the subcategory already has a parent.
func (c *Category) AddSubcategory(sub *Category) error {
	if sub.Parent != nil {
		return fmt.Errorf("cannot move subcategory %s to another parent: %s -> %s", sub.Name, sub.Parent.Name, c.Name)
	}
	c.Subcategories = append(c.Subcategories, sub)
	sub.Parent = c
	c.recomputeTotals()
	return nil
}

// AddTransaction assigns the given transaction to this Category, sets its
// Category pointer, and increments this Category's TotalAmount.
// It returns an error if the transaction is already assigned to a category.
func (c *Category) AddTransaction(txn *Transaction) error {
	if txn.Category != nil {
		return fmt.Errorf("cannot move transaction to another parent: %s -> %s", txn.Category.Name, c.Name)
	}
	c.Transactions = append(c.Transactions, txn)
	c.addToTotalAmount(txn.Amount)
	txn.Category = c
	return nil
}

// Categories groups the root‐level income, expense, and ignored categories
// for the financial application.
type Categories struct {
	// Income is the collection of revenue categories.
	Income *Category `json:"income"`
	// Expenses is the collection of spending categories.
	Expenses *Category `json:"expenses"`
	// Ignored holds categories excluded from reporting or processing.
	Ignored *Category `json:"ignored,omitempty"`
}

func NewCategories() *Categories {
	return &Categories{
		Income:   NewCategory("Income"),
		Expenses: NewCategory("Expenses"),
		Ignored:  NewCategory("Ignored"),
	}
}

// UnmarshalJSON implements custom JSON unmarshaling for Categories. After
// unmarshaling the raw data, it walks each category tree to restore Parent
// pointers and transaction Category links.
func (c *Categories) UnmarshalJSON(data []byte) error {
	// raw is an alias to avoid infinite recursion when unmarshaling
	type raw Categories
	var aux raw

	// Unmarshal into the alias type
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Restore parent/child relationships for each top‐level list
	rebuildTree := func(cat *Category) {
		if cat != nil {
			cat.rebuildTree()
		}
	}
	rebuildTree(aux.Income)
	rebuildTree(aux.Expenses)
	rebuildTree(aux.Ignored)

	// Assign back to the receiver
	*c = Categories(aux)
	return nil
}
