package gago

import (
	"os"
	"testing"
)

//TestReport Test antisampling and concurrency with batching
func TestReport(t *testing.T) {

	authFile := os.Getenv("GAGO_AUTH")
	analyticsreportingService, analyticsService := Authenticate(authFile)

	acc := GetAccountSummary(analyticsService)

	var req = GoogleAnalyticsRequest{
		Service:    analyticsreportingService,
		ViewID:     acc.viewID[0],
		Start:      "7DaysAgo",
		End:        "yesterday",
		Dimensions: "ga:date",
		Metrics:    "ga:sessions,ga:users",
		MaxRows:    100,
		AntiSample: false}

	report := GoogleAnalytics(req)

	if report.FetchedRowCount == 0 {
		t.Errorf("No rows fetched!")
	}
}
