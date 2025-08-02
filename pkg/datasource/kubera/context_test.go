package kubera

import (
	"context"
	"net/http"
	"os"
	"testing"

	clients "github.com/wyvernzora/personal-finance-mcp/internal/clients/kubera"
	kubera "github.com/wyvernzora/personal-finance-mcp/internal/clients/kubera"
)

// fakeClient implements the Kubera Client interface for testing.
type fakeClient struct{}

// GetPortfolio returns a predetermined portfolio for testing.
func (f *fakeClient) GetPortfolio(ctx context.Context) (*clients.Portfolio, error) {
	return &clients.Portfolio{
		Id:   "test-id",
		Name: "Test Portfolio",
		Assets: []*clients.AssetPosition{
			{
				Position: clients.Position{
					Id:      "a1",
					Name:    "Asset1",
					Ticker:  "AST1",
					Type:    "investment",
					Subtype: "stock",
					Value:   clients.Value{Amount: 100, Currency: "USD"},
				},
				Investable: "yes",
				Liquidity:  "high",
				AssetClass: "equity",
			},
		},
		Debts: []*clients.DebtPosition{
			{
				Position: clients.Position{
					Id:    "d1",
					Name:  "Debt1",
					Type:  "loan",
					Value: clients.Value{Amount: 40, Currency: "USD"},
				},
			},
		},
	}, nil
}

// TestInjectCredentialsFromEnvironment verifies that when all required
// environment variables are set, InjectCredentialsFromEnvironment returns
// an HTTPContextFunc that injects a non-nil Kubera Client into the context.
func TestInjectCredentialsFromEnvironment_Success(t *testing.T) {
	// Set required environment variables
	os.Setenv("KUBERA_API_KEY", "testKey")
	os.Setenv("KUBERA_API_SECRET", "testSecret")
	os.Setenv("KUBERA_PORTFOLIO_ID", "testPort")
	defer func() {
		os.Unsetenv("KUBERA_API_KEY")
		os.Unsetenv("KUBERA_API_SECRET")
		os.Unsetenv("KUBERA_PORTFOLIO_ID")
	}()

	// Obtain injector and apply to a background context
	injector := InjectCredentialsFromEnvironment()
	ctx := injector(context.Background(), &http.Request{})

	// FromContext should return a valid Kubera Client
	client := kubera.FromContext(ctx)
	if client == nil {
		t.Fatal("expected a non-nil Kubera Client in context")
	}
}

// TestInjectCredentialsFromEnvironment_MissingEnv verifies that if any
// required environment variable is missing, requireEnv will fatally exit.
// Because log.Fatalf calls os.Exit, we cannot catch it directly here.
// This test ensures that at least one missing var triggers a failure path
// by checking that Injector panics on missing configuration.
func TestInjectCredentialsFromEnvironment_MissingEnv(t *testing.T) {
	// Clear all relevant env vars
	os.Unsetenv("KUBERA_API_KEY")
	os.Unsetenv("KUBERA_API_SECRET")
	os.Unsetenv("KUBERA_PORTFOLIO_ID")

	// requireEnv should panic on missing vars
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when environment variable is missing")
		}
	}()
	_ = InjectCredentialsFromEnvironment()
}
