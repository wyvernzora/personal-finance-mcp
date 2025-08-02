package lunchmoney

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	lmapi "github.com/wyvernzora/personal-finance-mcp/internal/clients/lunch_money"
)

// InjectCredentialsFromEnvironment returns an HTTP context injector function that
// loads the LunchMoney API token from the environment variable "LUNCHMONEY_TOKEN".
// It returns a function that injects a LunchMoney client configured with the token
// into the context of incoming HTTP requests. Panics if the environment variable
// is not set or empty.
func InjectCredentialsFromEnvironment() func(ctx context.Context, req *http.Request) context.Context {
	token := requireEnv("LUNCHMONEY_TOKEN")
	return lmapi.WithLunchMoneyCredentials(token)
}

// requireEnv retrieves the value of the named environment variable.
// It panics if the variable is not set or is empty, logging the error message.
func requireEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		msg := fmt.Sprintf("environment variable %q must be set", name)
		log.Print(msg)
		panic(msg)
	}
	return val
}
