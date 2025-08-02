package types

// Transaction represents a financial transaction with a date, payee, and amount.
// Category is set when added to a Category tree, and user/system descriptions are available.
type Transaction struct {
	// Date is the date when the transaction occurred.
	Date Date `json:"date"`
	// Payee is the counterparty or recipient of the transaction.
	Payee string `json:"payee"`
	// Amount is the monetary value of the transaction.
	Amount Money `json:"amount"`
	// Description is an optional text provided by the end user for context.
	Description string `json:"description,omitempty"`
	// Category is a reference to the assigned Category; omitted from JSON.
	// It is automatically set by Category.AddTransaction.
	// To restore links after JSON unmarshaling, call rebuildTree on the root Category.
	Category *Category `json:"-"`
	// AnnotatedObject holds user and system descriptions for the transaction.
	AnnotatedObject
}

// NewTransaction creates a Transaction with the given date, payee, and amount.
// The returned Transaction has no Category (nil) and initialized Describable metadata.
func NewTransaction(date Date, payee string, amount Money) *Transaction {
	return &Transaction{
		Date:            date,
		Payee:           payee,
		Amount:          amount,
		AnnotatedObject: NewAnnotatedObject(),
	}
}
