package gago

import (
	"fmt"
	"os"
	"testing"
)

//TestReport Test antisampling and concurrency with batching
func TestAccounts(t *testing.T) {
	if os.Getenv("GAGO_AUTH") == "" {
		fmt.Println("Skip test, no auth")
		return
	}

	authFile := os.Getenv("GAGO_AUTH")

	_, analyticsService := Authenticate(authFile)

	GetAccountSummary(analyticsService)

}
