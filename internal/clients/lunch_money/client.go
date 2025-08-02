package lunchmoney

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// BASE_URL is the default Lunch Money API base URL.
const BASE_URL = "https://dev.lunchmoney.app"

type Client interface {
	ListTags(ctx context.Context) (Tags, error)
	ListCategories(ctx context.Context) (Categories, error)
	ListTransactions(ctx context.Context, startDate, endDate string) (Transactions, error)
}

// Client is a Lunch Money API client. It embeds an http.Client and holds auth and base URL config.
type client struct {
	*http.Client
	authToken string
	baseUrl   string
}

// NewClient creates and returns a new Client initialized with the provided auth token.
func NewClient(token string) Client {
	return &client{
		Client:    http.DefaultClient,
		authToken: token,
		baseUrl:   BASE_URL,
	}
}

// get sends an HTTP GET request to the client's base URL, appends the given path and query parameters,
// and returns the raw response body or an error.
func (c *client) get(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	u, _ := url.Parse(c.baseUrl)
	u.Path = path

	// Add query parameters
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return c.send(ctx, http.MethodGet, u.String(), nil)
}

// send constructs and executes an HTTP request with the specified method, URL, and body,
// sets the Authorization header using the client's token, verifies a 2xx status code,
// and returns the response body or an error.
func (c *client) send(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	// Execute
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status %d: %s", resp.StatusCode, string(body))
	}

	// Read and return body
	return io.ReadAll(resp.Body)
}

// Context key type and value
type lmClientKeyType struct{}

var lmClientKey = lmClientKeyType{}

// WithLunchMoneyCredentials returns an HTTP context function that initializes a new LunchMoney client using
// the supplied API token and stores the client in the context for future use.
func WithLunchMoneyCredentials(token string) func(ctx context.Context, r *http.Request) context.Context {
	return WithLunchMoneyClient(NewClient(token))
}

// WithLunchMoneyCredentials returns an HTTP context function that stores the supplied client in the context for future use.
func WithLunchMoneyClient(client Client) func(ctx context.Context, r *http.Request) context.Context {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, lmClientKey, client)
	}
}

// FromContext retrieves the LunchMoney API client stored in the context.
// It panics if no client is present.
func FromContext(ctx context.Context) Client {
	client, ok := ctx.Value(lmClientKey).(Client)
	if !ok {
		panic("Lunch Money client not found in context")
	}
	return client
}
