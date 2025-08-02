package lunchmoney

import (
	"context"
	"net/http"
	"os"
	"testing"

	lmapi "github.com/wyvernzora/personal-finance-mcp/internal/clients/lunch_money"
)

func TestInjectCredentialsFromEnvironment_Success(t *testing.T) {
	// Setup environment variable
	os.Setenv("LUNCHMONEY_TOKEN", "testtoken123")
	defer os.Unsetenv("LUNCHMONEY_TOKEN")

	injector := InjectCredentialsFromEnvironment()

	// Apply injector to context and dummy request
	ctx := context.Background()
	req := &http.Request{}
	ctxOut := injector(ctx, req)

	// Check if LunchMoney client is injected in context
	client := lmapi.FromContext(ctxOut)
	if client == nil {
		t.Fatal("expected LunchMoney client in context, got nil")
	}
}

func TestInjectCredentialsFromEnvironment_MissingToken(t *testing.T) {
	// Clear environment variable to simulate missing token
	os.Unsetenv("LUNCHMONEY_TOKEN")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic due to missing LUNCHMONEY_TOKEN environment variable")
		}
	}()

	_ = InjectCredentialsFromEnvironment()
}
