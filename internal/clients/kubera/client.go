// Package kubera provides a client for interacting with the Kubera API.
package kubera

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// BASE_URL is the default Kubera API base URL.
const BASE_URL = "https://api.kubera.com/api"

// Client defines the operations supported by the Kubera API client.
type Client interface {
	// GetPortfolio fetches the portfolio data for the configured portfolio ID.
	GetPortfolio(ctx context.Context) (*Portfolio, error)
}

// client implements the Client interface and handles authentication and request signing.
type client struct {
	*http.Client
	apiKey      string
	apiSecret   string
	portfolioId string
	baseUrl     string
}

// NewClient creates a new Kubera API client configured with apiKey, apiSecret, and portfolioId.
func NewClient(apiKey, apiSecret, portfolioId string) Client {
	return &client{
		Client:      http.DefaultClient,
		apiKey:      apiKey,
		apiSecret:   apiSecret,
		portfolioId: portfolioId,
		baseUrl:     BASE_URL,
	}
}

// get constructs and signs a GET request to the specified API path, executes it, and returns the response bytes.
func (c *client) get(ctx context.Context, path string) ([]byte, error) {
	method := http.MethodGet
	fullURL := c.baseUrl + path

	// Timestamp for signature
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Build signing string: apiKey + timestamp + method + path + body (empty)
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL for signing: %w", err)
	}
	pathToSign := u.RequestURI()
	signingString := c.apiKey + timestamp + method + pathToSign

	// Compute HMAC-SHA256 signature
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(signingString))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set auth headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-token", c.apiKey)
	req.Header.Set("x-timestamp", timestamp)
	req.Header.Set("x-signature", signature)

	// Perform request
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read and return response body
	return io.ReadAll(resp.Body)
}

type kuberaClientKeyType struct{}

var kuberaClientKey = kuberaClientKeyType{}

// WithKuberaCredentials returns an HTTPContextFunc that initializes a Kubera Client using the provided credentials
// and stores that client in the context for future use.
func WithKuberaCredentials(apiKey, apiSecret, portfolioId string) func(ctx context.Context, r *http.Request) context.Context {
	client := NewClient(apiKey, apiSecret, portfolioId)
	return WithKuberaClient(client)
}

// WithKuberaClient returns an HTTPContextFunc that stores the supplied Kubera client in the context.
func WithKuberaClient(client Client) func(ctx context.Context, r *http.Request) context.Context {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, kuberaClientKey, client)
	}
}

// FromContext fetches the Kubera Client from the context. It panics if no client is found.
func FromContext(ctx context.Context) Client {
	client, ok := ctx.Value(kuberaClientKey).(Client)
	if !ok {
		panic("Kubera client not found in context")
	}
	return client
}
