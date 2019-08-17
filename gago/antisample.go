package gago

import (
	"fmt"

	ga "google.golang.org/api/analyticsreporting/v4"
)

func makeAntiSampleRequestList(gagoRequest *GoogleAnalyticsRequest) [][]*ga.ReportRequest {
	fmt.Println("antisampling")
	// do call to test if report is sampled
	test := GoogleAnalyticsRequest{
		Service:    gagoRequest.Service,
		ViewID:     gagoRequest.ViewID,
		Start:      gagoRequest.Start,
		End:        gagoRequest.End,
		Dimensions: gagoRequest.Dimensions,
		Metrics:    gagoRequest.Metrics,
		MaxRows:    18}

	testResponse := GoogleAnalytics(test)

	if testResponse.SamplesReadCounts == nil ||
		testResponse.SamplingSpaceSizes == nil {
		//if not, return normal list
		fmt.Println("No sampling found")
		return makeRequestList(gagoRequest)
	}

	fmt.Println(testResponse)

	//fmt.Println("sampling found: ", (testResponse.SamplesReadCounts[0]/testResponse.SamplingSpaceSizes[0])*100)

	// if sampled, fetch exploratory sessions call
	var explore = GoogleAnalyticsRequest{
		Service:    gagoRequest.Service,
		ViewID:     gagoRequest.ViewID,
		Start:      gagoRequest.Start,
		End:        gagoRequest.End,
		Dimensions: "ga:date",
		Metrics:    "ga:sessions",
		MaxRows:    9999}
	exploreResponse := GoogleAnalytics(explore)
	fmt.Println("Explore found", exploreResponse.Totals)

	// work out date ranges to fetch

	// construct new GoogleAnalyticsRequest objects via makeRequestList(gagoRequest)

	// create new ga.ReportRequests

	// return
	return makeRequestList(gagoRequest)

}
