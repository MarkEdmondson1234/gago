package gago

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	ga "google.golang.org/api/analyticsreporting/v4"
)

func makeRequestList(gagoRequest *GoogleAnalyticsRequest) [][]*ga.ReportRequest {

	gagoRequest.maxPages = (gagoRequest.MaxRows / (gagoRequest.PageLimit * apiBatchLimit)) + 1
	fmt.Println("maxPages: ", gagoRequest.maxPages)

	requestList := make([][]*ga.ReportRequest, gagoRequest.maxPages)
	fetchMore := true
	gagoRequest.fetchedRows = 0
	for i := 0; fetchMore; i++ {
		//fmt.Println("paging: ", i, fetchMore)

		// a loop around 5 requests
		reqp := make([]*ga.ReportRequest, apiBatchLimit)
		for j := range reqp {
			//fmt.Println("ps", gagoRequest.pageSize, " pt", gagoRequest.pageToken, " pl",
			//	gagoRequest.PageLimit, " mr", gagoRequest.MaxRows, "fr", gagoRequest.fetchedRows)
			req := makeRequest(*gagoRequest)
			reqp[j] = req

			gagoRequest.pageToken = strconv.FormatInt(gagoRequest.fetchedRows+gagoRequest.pageSize, 10)
			gagoRequest.fetchedRows = gagoRequest.fetchedRows + gagoRequest.pageSize

			// stop fetching if we adjusted the pageSize down
			// stop fetching if we've reached maxRows
			if gagoRequest.pageSize < gagoRequest.PageLimit ||
				(gagoRequest.MaxRows > 0 && gagoRequest.fetchedRows >= gagoRequest.MaxRows) {
				fmt.Println("dont fetchmore")
				fetchMore = false
				break
			}

			// do we need to modify pageSize for next loop?
			if gagoRequest.MaxRows > 0 && (gagoRequest.fetchedRows+gagoRequest.pageSize) > gagoRequest.MaxRows {
				gagoRequest.pageSize = gagoRequest.MaxRows - gagoRequest.fetchedRows + 1
			}
		}

		requestList[i] = reqp

	}

	return requestList
}

func fetchConcurrentReport(requestList [][]*ga.ReportRequest, gagoRequest GoogleAnalyticsRequest) []*ga.GetReportsResponse {
	// 10 concurrent requests per view (profile) (cannot be increased)
	// 1- 1-10, 2 - 11-20 etc.
	responses := make([]*ga.GetReportsResponse, gagoRequest.maxPages)

	fmt.Println("maxPages: ", gagoRequest.maxPages)

	//fmt.Println("requestList>", requestList)

	responseIndex := 0
	for i := 0; i < len(requestList); i += concurrencyLimit {

		//fmt.Println("batch: ", i, min(i+concurrencyLimit, len(requestList)))
		batch := requestList[i:min(i+concurrencyLimit, len(requestList))]
		var wg sync.WaitGroup
		fmt.Println("api concurrency size:", len(batch))

		wg.Add(len(batch))

		for j, request := range batch {
			//fmt.Println("j", j, "request", request)
			// fetch requests
			go func(j int, request []*ga.ReportRequest, gagoRequest GoogleAnalyticsRequest, responseIndex int) {
				defer wg.Done()
				responses[responseIndex] = fetchReport(gagoRequest, request)
				//fmt.Println("responseIndex: ", responseIndex)
			}(j, request, gagoRequest, responseIndex)
			responseIndex++
		}

		wg.Wait()

	}

	// fmt.Println("responseIndex", responseIndex)
	// for _, res := range responses {
	// 	//fmt.Println("response i>", i)
	// 	if res == nil {
	// 		fmt.Println("res=nil")
	// 		continue
	// 	}
	// 	for _, report := range res.Reports {
	// 		fmt.Println(report.Data.RowCount)
	// 	}
	// }

	return responses
}

//makeRequest creates the request(s) for fetchReport
// start and end are YYYY-mm-dd
// dimensions and metrics are ga:dim1,ga:dim2 and ga:metric1,ga:metric2
func makeRequest(gagoRequest GoogleAnalyticsRequest) *ga.ReportRequest {

	// slice of length 1 of type *ga.DateRange
	daterangep := make([]*ga.DateRange, 1)
	// Fill the 1st element with a pointer to a ga.DateRange
	daterangep[0] = &ga.DateRange{StartDate: gagoRequest.Start, EndDate: gagoRequest.End}

	dimSplit := strings.Split(gagoRequest.Dimensions, ",")
	dimp := make([]*ga.Dimension, len(dimSplit))
	for i, dim := range dimSplit {
		dimp[i] = &ga.Dimension{Name: dim}
	}

	metSplit := strings.Split(gagoRequest.Metrics, ",")
	metp := make([]*ga.Metric, len(metSplit))
	for i, met := range metSplit {
		metp[i] = &ga.Metric{Expression: met}
	}

	requests := ga.ReportRequest{}
	requests.DateRanges = daterangep
	requests.Dimensions = dimp
	requests.Metrics = metp
	requests.IncludeEmptyRows = true
	requests.PageSize = gagoRequest.pageSize
	requests.ViewId = gagoRequest.ViewID
	requests.SamplingLevel = "LARGE"
	requests.PageToken = gagoRequest.pageToken

	// print out json request
	js, _ := requests.MarshalJSON()
	fmt.Println("\n# Request:", string(js))

	return &requests
}

// FetchReport Perform the GAv4 API request
func fetchReport(
	gagoRequest GoogleAnalyticsRequest,
	reports []*ga.ReportRequest) *ga.GetReportsResponse {

	reportreq := &ga.GetReportsRequest{ReportRequests: reports, UseResourceQuotas: gagoRequest.UseResourceQuotas}

	// Max 5 reportreq's at once
	report, err := gagoRequest.Service.Reports.BatchGet(reportreq).Do()
	if err != nil {
		log.Fatal(err)
	}

	if report == nil {
		log.Fatal("Nil report:", reportreq)
	}

	//js, _ := json.Marshal(report)
	//fmt.Println("\n## fetched ", string(js))

	return report

}
