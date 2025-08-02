package lunchmoney

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestListTags_Success(t *testing.T) {
	// Prepare sample tags JSON
	tags := []*Tag{
		{Id: 101, Name: "food", Description: "Food expenses", IsArchived: false},
		{Id: 202, Name: "travel", Description: "Travel costs", IsArchived: true},
	}
	payload, err := json.Marshal(tags)
	if err != nil {
		t.Fatalf("marshal sample tags: %v", err)
	}

	cli := newTestClient("tok", func(req *http.Request) (*http.Response, error) {
		// verify correct endpoint
		want := BASE_URL + "/v1/tags"
		if req.URL.String() != want {
			t.Errorf("URL = %q; want %q", req.URL.String(), want)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(string(payload))),
			Header:     make(http.Header),
		}, nil
	})

	out, err := cli.ListTags(context.Background())
	if err != nil {
		t.Fatalf("ListTags returned error: %v", err)
	}
	if len(out) != len(tags) {
		t.Fatalf("len(out) = %d; want %d", len(out), len(tags))
	}
	for _, tag := range tags {
		got, ok := out[tag.Id]
		if !ok {
			t.Errorf("missing tag ID %d in result", tag.Id)
			continue
		}
		if got.Name != tag.Name || got.Description != tag.Description || got.IsArchived != tag.IsArchived {
			t.Errorf("got %+v; want %+v", got, tag)
		}
	}
}

func TestListTags_HTTPError(t *testing.T) {
	cli := newTestClient("tok", func(_ *http.Request) (*http.Response, error) {
		return nil, errors.New("network failure")
	})
	_, err := cli.ListTags(context.Background())
	if err == nil {
		t.Fatal("expected error on network failure, got nil")
	}
	if !strings.Contains(err.Error(), "failed to call Lunch Money API") {
		t.Errorf("error = %q; want it to contain %q", err.Error(), "failed to call Lunch Money API")
	}
}

func TestListTags_BadJSON(t *testing.T) {
	cli := newTestClient("tok", func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("not-json")),
			Header:     make(http.Header),
		}, nil
	})
	_, err := cli.ListTags(context.Background())
	if err == nil {
		t.Fatal("expected JSON unmarshal error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to deserialize response") {
		t.Errorf("error = %q; want it to contain %q", err.Error(), "failed to deserialize response")
	}
}
