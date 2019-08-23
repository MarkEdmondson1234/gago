package gago

import (
	"os"
	"testing"
)

//TestReport Test antisampling and concurrency with batching
func TestAccounts(t *testing.T) {
	if os.Getenv("GAGO_AUTH") == "" {
		t.Skip("Skip test, no auth")
	}

	authFile := os.Getenv("GAGO_AUTH")

	_, analyticsService := Authenticate(authFile)

	GetAccountSummary(analyticsService)

}
