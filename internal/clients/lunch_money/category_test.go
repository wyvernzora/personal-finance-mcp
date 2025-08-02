package lunchmoney

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestListCategories_Success(t *testing.T) {
	orig := http.DefaultClient
	defer func() { http.DefaultClient = orig }()

	const jsonBody = `{
		"categories": [
			{
				"id": 1,
				"name": "Food",
				"description": "All food expenses",
				"order": 10,
				"is_income": false,
				"exclude_from_budget": false,
				"exclude_from_totals": false,
				"is_archived": false,
				"archived_on": "",
				"updated_at": "2023-01-01T00:00:00Z",
				"created_at": "2022-01-01T00:00:00Z"
			}
		]
	}`

	http.DefaultClient = &http.Client{
		Transport: &fakeTransport{fn: func(req *http.Request) (*http.Response, error) {
			wantURL := BASE_URL + "/v1/categories?format=nested"
			if req.URL.String() != wantURL {
				t.Errorf("request URL = %q; want %q", req.URL.String(), wantURL)
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(jsonBody)),
				Header:     make(http.Header),
			}, nil
		}},
	}

	client := NewClient("dummy-token")
	cats, err := client.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("ListCategories returned error: %v", err)
	}
	if len(cats) != 1 {
		t.Fatalf("len(categories) = %d; want 1", len(cats))
	}
	cat, ok := cats[1]
	if !ok {
		t.Fatalf("expected category ID 1 in map")
	}
	if cat.Id != 1 {
		t.Errorf("Id = %d; want 1", cat.Id)
	}
	if cat.Name != "Food" {
		t.Errorf("Name = %q; want %q", cat.Name, "Food")
	}
	if cat.Description != "All food expenses" {
		t.Errorf("Description = %q; want %q", cat.Description, "All food expenses")
	}
}

func TestListCategories_HTTPError(t *testing.T) {
	orig := http.DefaultClient
	defer func() { http.DefaultClient = orig }()

	http.DefaultClient = &http.Client{
		Transport: &fakeTransport{fn: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network failure")
		}},
	}

	client := NewClient("token")
	_, err := client.ListCategories(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to call Lunch Money API") {
		t.Errorf("error = %q; want it to contain %q", err.Error(), "failed to call Lunch Money API")
	}
}

func TestListCategories_BadStatus(t *testing.T) {
	orig := http.DefaultClient
	defer func() { http.DefaultClient = orig }()

	http.DefaultClient = &http.Client{
		Transport: &fakeTransport{fn: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("server error")),
				Header:     make(http.Header),
			}, nil
		}},
	}

	client := NewClient("token")
	_, err := client.ListCategories(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to call Lunch Money API") {
		t.Errorf("error = %q; want it to contain %q", err.Error(), "failed to call Lunch Money API")
	}
}

func TestListCategories_BadJSON(t *testing.T) {
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

	client := NewClient("token")
	_, err := client.ListCategories(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to deserialize response") {
		t.Errorf("error = %q; want it to contain %q", err.Error(), "failed to deserialize response")
	}
}
