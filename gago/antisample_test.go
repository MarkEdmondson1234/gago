package gago

import (
	"os"
	"testing"
)

//TestReport Test antisampling and concurrency with batching
func TestAntisample(t *testing.T) {
	if os.Getenv("GAGO_AUTH") == "" {
		t.Skip("Skip test, no auth")
	}

	authFile := os.Getenv("GAGO_AUTH")
	analyticsreportingService, _ := Authenticate(authFile)

	var req = GoogleAnalyticsRequest{
		Service:    analyticsreportingService,
		ViewID:     "106249469",
		Start:      "2016-07-01",
		End:        "2019-08-01",
		Dimensions: "ga:date,ga:sourceMedium,ga:landingPagePath,ga:source,ga:hour,ga:minute,ga:eventCategory",
		Metrics:    "ga:sessions,ga:users",
		MaxRows:    -1,
		AntiSample: true}

	report := GoogleAnalytics(req)

	if report.FetchedRowCount == 0 {
		t.Errorf("No rows fetched!")
	}
}
