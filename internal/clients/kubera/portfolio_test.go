package kubera

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetPortfolio_Success(t *testing.T) {
	// Prepare a valid API payload with errorCode=0
	want := &Portfolio{
		Id:   "port99",
		Name: "TestPort",
		Assets: []*AssetPosition{
			{Position: Position{Id: "a1", Name: "Asset1", Value: Value{Amount: 100, Currency: "USD"}}},
		},
		Debts: []*DebtPosition{
			{Position: Position{Id: "d1", Name: "Debt1", Value: Value{Amount: 40, Currency: "USD"}}},
		},
	}
	payload := struct {
		Data      *Portfolio `json:"data"`
		ErrorCode int64      `json:"errorCode"`
	}{Data: want, ErrorCode: 0}
	body, _ := json.Marshal(payload)

	cli := newTestClient("k", "s", "port99", func(req *http.Request) (*http.Response, error) {
		// verify request path
		expectURL := BASE_URL + "/v3/data/portfolio/port99"
		if req.URL.String() != expectURL {
			t.Errorf("URL = %q; want %q", req.URL.String(), expectURL)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Header:     make(http.Header),
		}, nil
	})

	got, err := cli.GetPortfolio(context.Background())
	if err != nil {
		t.Fatalf("GetPortfolio error: %v", err)
	}
	if got.Id != want.Id || got.Name != want.Name {
		t.Errorf("got %+v; want %+v", got, want)
	}
	if len(got.Assets) != 1 || got.Assets[0].Id != "a1" {
		t.Errorf("assets = %+v; want one asset with Id a1", got.Assets)
	}
	if len(got.Debts) != 1 || got.Debts[0].Id != "d1" {
		t.Errorf("debts = %+v; want one debt with Id d1", got.Debts)
	}
}

func TestGetPortfolio_APIError(t *testing.T) {
	// Simulate errorCode > 0
	resp := `{"data":null,"errorCode":123}`
	cli := newTestClient("k", "s", "pid", func(_ *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(resp)), Header: make(http.Header)}, nil
	})

	_, err := cli.GetPortfolio(context.Background())
	if err == nil || !strings.Contains(err.Error(), "Kubera API error: 123") {
		t.Errorf("err = %v; want error containing Kubera API error: 123", err)
	}
}

func TestGetPortfolio_NetworkFailure(t *testing.T) {
	cli := newTestClient("k", "s", "pid", func(_ *http.Request) (*http.Response, error) {
		return nil, errors.New("conn refused")
	})

	_, err := cli.GetPortfolio(context.Background())
	if err == nil || !strings.Contains(err.Error(), "failed to call Kubera API") {
		t.Errorf("err = %v; want error prefix failed to call Kubera API", err)
	}
}

func TestGetPortfolio_DeserializeError(t *testing.T) {
	cli := newTestClient("k", "s", "pid", func(_ *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("not json")), Header: make(http.Header)}, nil
	})

	_, err := cli.GetPortfolio(context.Background())
	if err == nil || !strings.Contains(err.Error(), "failed to deserialize response") {
		t.Errorf("err = %v; want error containing failed to deserialize response", err)
	}
}
