package kubera

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// fakeTransport lets us stub HTTP responses for client.get.
type fakeTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (f *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.fn(req)
}

// newTestClient constructs a *client whose HTTP transport is overridden by fn.
func newTestClient(apiKey, apiSecret, portfolioId string, fn func(req *http.Request) (*http.Response, error)) *client {
	cli := NewClient(apiKey, apiSecret, portfolioId).(*client)
	cli.Client = &http.Client{Transport: &fakeTransport{fn: fn}}
	return cli
}

func TestClientGet_Success(t *testing.T) {
	const body = "hello world"
	// Create a test client that returns 200 + body for path "/foo"
	cli := newTestClient("k", "s", "pid", func(req *http.Request) (*http.Response, error) {
		want := BASE_URL + "/foo"
		if req.URL.String() != want {
			t.Errorf("URL = %q; want %q", req.URL.String(), want)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	data, err := cli.get(context.Background(), "/foo")
	if err != nil {
		t.Fatalf("get returned error: %v", err)
	}
	if string(data) != body {
		t.Errorf("body = %q; want %q", string(data), body)
	}
}

func TestClientGet_NetworkError(t *testing.T) {
	cli := newTestClient("k", "s", "pid", func(req *http.Request) (*http.Response, error) {
		return nil, io.ErrUnexpectedEOF
	})
	_, err := cli.get(context.Background(), "/foo")
	if err == nil || !strings.Contains(err.Error(), "request failed") {
		t.Errorf("err = %v; want network error containing 'request failed'", err)
	}
}

func TestClientGet_BadStatus(t *testing.T) {
	cli := newTestClient("k", "s", "pid", func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(strings.NewReader("oops")),
			Header:     make(http.Header),
		}, nil
	})
	_, err := cli.get(context.Background(), "/foo")
	if err == nil || !strings.Contains(err.Error(), "bad status 500") {
		t.Errorf("err = %v; want error containing 'bad status 500'", err)
	}
}

func TestWithKuberaClientAndFromContext(t *testing.T) {
	ctx := context.Background()
	req, _ := http.NewRequest("GET", "/", nil)
	fn := WithKuberaCredentials("key", "secret", "pid")
	newCtx := fn(ctx, req)
	c := FromContext(newCtx)
	if c == nil {
		t.Fatal("FromContext returned nil, expected non-nil client")
	}
}

func TestFromContext_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when no client in context")
		}
	}()
	_ = FromContext(context.Background())
}
