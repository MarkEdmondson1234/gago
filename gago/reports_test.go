package gago

import (
	"os"
	"testing"
)

//TestReport Test antisampling and concurrency with batching
func TestReport(t *testing.T) {

	authFile := os.Getenv("GAGO_AUTH")
	analyticsreportingService, _ := Authenticate(authFile)

	var req = GoogleAnalyticsRequest{
		Service:    analyticsreportingService,
		ViewID:     "106249469",
		Start:      "2016-07-01",
		End:        "2019-08-01",
		Dimensions: "ga:date,ga:sourceMedium,ga:landingPagePath,ga:source,ga:hour,ga:minute,ga:eventCategory",
		Metrics:    "ga:sessions,ga:users",
		MaxRows:    0,
		AntiSample: true}

	report := GoogleAnalytics(req)

	if report.Totals[0] != "11665" {
		t.Errorf("Expected report.Totals = '11665' but got %s", report.Totals[0])
	}
}

//TestAccountSummary Test account summary
func TestAccountSummary(t *testing.T) {

	authFile := os.Getenv("GAGO_AUTH")
	_, analyticsService := Authenticate(authFile)

	GetAccountSummary(analyticsService)
	//Output:
	//accountId accountName webPropertyId webPropertyName viewId viewName
	//47480439 MarkEdmondson UA-47480439-2 MarkEdmondson oldGA 81416156 Live Blog
	//47480439 MarkEdmondson UA-47480439-1 markedmondson.me 81416941 All Web Site Data
	//54019251 Sunholo Websites UA-54019251-4 Mark's GA360 WebProperty 106249469 WF GA Premium E-commerce
}
