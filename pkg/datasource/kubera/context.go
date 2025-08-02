package kubera

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wyvernzora/personal-finance-mcp/internal/clients/kubera"
)

// InjectCredentialsFromEnvironment loads Kubera API credentials (API key, secret, and portfolio ID)
// from environment variables and returns an HTTPContextFunc that injects the configured Kubera client into the request context.
func InjectCredentialsFromEnvironment() func(ctx context.Context, req *http.Request) context.Context {
	apiKey := requireEnv("KUBERA_API_KEY")
	apiSecret := requireEnv("KUBERA_API_SECRET")
	portfolioId := requireEnv("KUBERA_PORTFOLIO_ID")

	return kubera.WithKuberaCredentials(apiKey, apiSecret, portfolioId)
}

// requireEnv retrieves the value of the named environment variable.
// It panics if the variable is not set or empty, enforcing required configuration.
// Used by InjectCredentialsFromEnvironment to ensure all credentials are present.
func requireEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		msg := fmt.Sprintf("environment variable %q must be set", name)
		log.Print(msg)
		panic(msg)
	}
	return val
}
