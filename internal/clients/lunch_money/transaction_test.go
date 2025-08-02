package lunchmoney

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/wyvernzora/personal-finance-mcp/pkg/types"
)

func TestListTransactions_Success(t *testing.T) {
	orig := http.DefaultClient
	defer func() { http.DefaultClient = orig }()

	const bodyJSON = `{
		"transactions": [
			{"id": 1, "date": "2023-01-01", "to_base": 12.34, "payee": "Alice", "original_name": "Alice",
			 "category_name": "Food", "category_group_name": "Expenses", "is_income": false,
			 "exclude_from_budget": false, "exclude_from_totals": false, "display_notes": "note"}
		]
	}`

	http.DefaultClient = &http.Client{
		Transport: &fakeTransport{fn: func(req *http.Request) (*http.Response, error) {
			// verify path and query parameters
			expectedPath := "/v1/transactions"
			if req.URL.Path != expectedPath {
				t.Errorf("unexpected path: got %q, want %q", req.URL.Path, expectedPath)
			}
			q := req.URL.Query()
			if q.Get("start_date") != "2023-01-01" {
				t.Errorf("start_date = %q, want %q", q.Get("start_date"), "2023-01-01")
			}
			if q.Get("end_date") != "2023-01-31" {
				t.Errorf("end_date = %q, want %q", q.Get("end_date"), "2023-01-31")
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(bodyJSON)),
				Header:     make(http.Header),
			}, nil
		}},
	}

	client := NewClient("token123")
	txs, err := client.ListTransactions(context.Background(), "2023-01-01", "2023-01-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("got %d transactions; want 1", len(txs))
	}
	tx := txs[0]
	if tx.Id != 1 {
		t.Errorf("Id = %d; want %d", tx.Id, 1)
	}
	if tx.Payee != "Alice" {
		t.Errorf("Payee = %q; want %q", tx.Payee, "Alice")
	}
	if tx.Amount != types.Money(123400) {
		t.Errorf("Amount = %d; want %d", tx.Amount, types.Money(123400))
	}
}

func TestListTransactions_HTTPError(t *testing.T) {
	orig := http.DefaultClient
	defer func() { http.DefaultClient = orig }()

	http.DefaultClient = &http.Client{
		Transport: &fakeTransport{fn: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network fail")
		}},
	}

	client := NewClient("tk")
	_, err := client.ListTransactions(context.Background(), "a", "b")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to call Lunch Money API") {
		t.Errorf("error = %q; want prefix %q", err.Error(), "failed to call Lunch Money API")
	}
}

func TestListTransactions_BadStatus(t *testing.T) {
	orig := http.DefaultClient
	defer func() { http.DefaultClient = orig }()

	http.DefaultClient = &http.Client{
		Transport: &fakeTransport{fn: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("err")),
				Header:     make(http.Header),
			}, nil
		}},
	}

	client := NewClient("tk")
	_, err := client.ListTransactions(context.Background(), "start", "end")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to call Lunch Money API") {
		t.Errorf("error = %q; want prefix %q", err.Error(), "failed to call Lunch Money API")
	}
}

func TestListTransactions_BadJSON(t *testing.T) {
	orig := http.DefaultClient
	defer func() { http.DefaultClient = orig }()

	http.DefaultClient = &http.Client{
		Transport: &fakeTransport{fn: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("not json")),
				Header:     make(http.Header),
			}, nil
		}},
	}

	client := NewClient("tk")
	_, err := client.ListTransactions(context.Background(), "s", "e")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to deserialize response") {
		t.Errorf("error = %q; want prefix %q", err.Error(), "failed to deserialize response")
	}
}
