package gago

import (
	"fmt"
	"strconv"

	"google.golang.org/api/analyticsreporting/v4"
	ga "google.golang.org/api/analyticsreporting/v4"
)

const antiSampleBatchSize = 250000

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
		MaxRows:    100}

	testResponse := GoogleAnalytics(test)

	if testResponse.SamplesReadCounts == nil ||
		testResponse.SamplingSpaceSizes == nil {
		//if not, return normal list
		fmt.Println("No sampling found")
		return makeRequestList(gagoRequest)
	}

	// update maxPages in request
	gagoRequest.maxPages = 1000000/gagoRequest.PageLimit + 1
	gagoRequest.fetchedRows = testResponse.RowCount + 1

	readCounts := float64(testResponse.SamplesReadCounts[0])
	samplingSize := float64(testResponse.SamplingSpaceSizes[0])

	fmt.Println("sampling found: ", (readCounts/samplingSize)*100)

	// if sampled, fetch exploratory sessions call
	var explore = GoogleAnalyticsRequest{
		Service:    gagoRequest.Service,
		ViewID:     gagoRequest.ViewID,
		Start:      gagoRequest.Start,
		End:        gagoRequest.End,
		Dimensions: "ga:year,ga:month,ga:day",
		Metrics:    "ga:sessions",
		MaxRows:    9999}
	exploreResponse := GoogleAnalytics(explore)
	fmt.Println("Explore found", exploreResponse.Totals)

	// work out date ranges to fetch
	sessionsSoFar := 0
	newStartDates := make([]string, len(exploreResponse.Rows))
	newEndDates := make([]string, len(exploreResponse.Rows))
	newRequestIndex := 1
	var previousDate string
	for i, row := range exploreResponse.Rows {

		if row == nil {
			break // hides some sins? why is there a nil row?
		}

		thisDate := join(row.Dimensions[0], "-",
			row.Dimensions[1], "-",
			row.Dimensions[2])

		if i == 0 {
			newStartDates[0] = thisDate
			continue
		}

		rowSession, _ := strconv.Atoi(row.Metrics[0])
		sessionsSoFar += rowSession
		if sessionsSoFar >= antiSampleBatchSize {
			newStartDates[newRequestIndex] = thisDate
			newEndDates[newRequestIndex-1] = previousDate
			newRequestIndex++
			sessionsSoFar = 0
		}

		previousDate = thisDate
	}
	newEndDates[newRequestIndex] = gagoRequest.End

	newStartDates = deleteEmptyStringSlice(newStartDates)
	newEndDates = deleteEmptyStringSlice(newEndDates)

	fmt.Println("start dates", newStartDates)
	fmt.Println("end dates", newEndDates)

	// construct new GoogleAnalyticsRequest objects via makeRequestList(gagoRequest)
	antiSampleRequests := make([][][]*analyticsreporting.ReportRequest, len(newStartDates))
	totalRequests := 0
	for i, date := range newStartDates {
		if date == "" {
			break
		}

		req := &GoogleAnalyticsRequest{
			Service:    gagoRequest.Service,
			ViewID:     gagoRequest.ViewID,
			Start:      newStartDates[i],
			End:        newEndDates[i],
			Dimensions: gagoRequest.Dimensions,
			Metrics:    gagoRequest.Metrics,
			MaxRows:    999999, // antisampling always gets all rows
			PageLimit:  gagoRequest.PageLimit,
			maxPages:   gagoRequest.maxPages,
			pageSize:   gagoRequest.PageLimit,
		}

		// create new ga.ReportRequests
		antiSampleRequests[i] = makeRequestList(req)

		totalRequests += len(antiSampleRequests[i])

	}

	//fmt.Println("antiSampleRequests>", antiSampleRequests)

	outputList := make([][]*analyticsreporting.ReportRequest, totalRequests)

	//remove one level of nesting
	tr := 0
	for _, ll := range antiSampleRequests {
		for _, lll := range ll {
			outputList[tr] = lll
			tr++
		}
	}

	fmt.Println("total requests: ", totalRequests, " outputList: ", len(outputList), "tr: ", tr)

	// return
	return outputList

}
