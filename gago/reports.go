package gago

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	ga "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/googleapi"
)

//GoogleAnalyticsRequest Make a request object to pass to GoogleAnalytics
type GoogleAnalyticsRequest struct {
	Service                                 *ga.Service
	ViewID, Start, End, Dimensions, Metrics string
	MaxRows, PageLimit                      int64
	UseResourceQuotas, AntiSample           bool
	pageSize, fetchedRows                   int64
	pageToken                               string
}

//GoogleAnalytics Make a request to the v4 Reporting API
func GoogleAnalytics(gagoRequest GoogleAnalyticsRequest) *ParseReport {

	// init ""
	if gagoRequest.MaxRows == 0 {
		gagoRequest.MaxRows = 1000
	}
	if gagoRequest.PageLimit == 0 {
		gagoRequest.PageLimit = 2
	}

	gagoRequest.pageSize = gagoRequest.PageLimit
	gagoRequest.MaxRows = gagoRequest.MaxRows - 1 //0 index

	if gagoRequest.MaxRows < gagoRequest.PageLimit {
		// if first page needs to fetch less than 10k default
		gagoRequest.pageSize = gagoRequest.MaxRows
	}

	maxPages := (gagoRequest.MaxRows / (gagoRequest.PageLimit * 5)) + 1

	responses := make([]*ga.GetReportsResponse, maxPages)

	fmt.Println("maxPages: ", maxPages)

	gagoRequest.fetchedRows = 0
	fetchMore := true

	requestList := make([][]*ga.ReportRequest, maxPages)

	for i := 0; fetchMore; i++ {
		//fmt.Println("paging: ", i, fetchMore)

		// a loop around 5 requests
		reqp := make([]*ga.ReportRequest, 5)
		for j := range reqp {
			// fmt.Println("ps", gagoRequest.pageSize, " pt", gagoRequest.pageToken, " pl", gagoRequest.PageLimit, " mr", gagoRequest.MaxRows, "fr", fetchedRows)
			req := makeRequest(gagoRequest)
			reqp[j] = req

			gagoRequest.pageToken = strconv.FormatInt(gagoRequest.fetchedRows+gagoRequest.pageSize, 10)
			gagoRequest.fetchedRows = gagoRequest.fetchedRows + gagoRequest.pageSize

			// stop fetching if we adjusted the pageSize down
			// stop fetching if we've reached maxRows
			if gagoRequest.pageSize < gagoRequest.PageLimit ||
				(gagoRequest.MaxRows > 0 && gagoRequest.fetchedRows >= gagoRequest.MaxRows) {
				// fmt.Println("dont fetchmore")
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

	// 10 concurrent requests per view (profile) (cannot be increased)
	// 1- 1-10, 2 - 11-20 etc.
	concurrencyLimit := 10
	concurrentRequests := ((len(requestList) - 1) / concurrencyLimit) + 1
	fmt.Println("concurrency: ", concurrentRequests)

	responseIndex := 0
	for i := 0; i < len(requestList); i += concurrencyLimit {

		batch := requestList[i:min(i+concurrencyLimit, len(requestList))]
		var wg sync.WaitGroup
		fmt.Println("batch size:", len(batch))

		wg.Add(len(batch))

		for j, request := range batch {
			// fetch requests
			go func(j int, request []*ga.ReportRequest, gagoRequest GoogleAnalyticsRequest) {
				defer wg.Done()
				responses[responseIndex] = fetchReport(gagoRequest, request)
				responseIndex++
			}(j, request, gagoRequest)

		}

		wg.Wait()

	}

	js, _ := json.MarshalIndent(responses, "", " ")
	fmt.Println("\n# All Responses:", string(js))

	parseReports, _ := ParseReportsResponse(responses, gagoRequest.fetchedRows)

	return parseReports

}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

//func worker(service *ga.Service)

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

	fmt.Println("fetching: pt", gagoRequest.pageToken,
		"ps:", gagoRequest.pageSize,
		"fr:", gagoRequest.fetchedRows)

	reportreq := &ga.GetReportsRequest{ReportRequests: reports, UseResourceQuotas: gagoRequest.UseResourceQuotas}

	// Max 5 reportreq's at once
	report, err := gagoRequest.Service.Reports.BatchGet(reportreq).Do()
	if err != nil {
		log.Fatal(err)
	}

	//js, _ := json.Marshal(report)
	//fmt.Println("\n## fetched ", string(js))

	return report

}

// ParseReportRow A parsed row of ParseReport
type ParseReportRow struct {
	Dimensions []string `json:"dimensions,omitempty"`
	Metrics    []string `json:"metrics,omitempty"`
}

// ParseReport A parsed Report after all batching and paging
type ParseReport struct {
	ColumnHeaderDimension []string                `json:"dimensionHeaderEntries,omitempty"`
	ColumnHeaderMetrics   []*ga.MetricHeaderEntry `json:"metricHeaderEntries,omitempty"`
	Rows                  []*ParseReportRow       `json:"values,omitempty"`
	DataLastRefreshed     string                  `json:"dataLastRefreshed,omitempty"`
	IsDataGolden          bool                    `json:"isDataGolden,omitempty"`
	Maximums              []string                `json:"maximums,omitempty"`
	Minimums              []string                `json:"minimums,omitempty"`
	RowCount              int64                   `json:"rowCount,omitempty"`
	SamplesReadCounts     googleapi.Int64s        `json:"samplesReadCounts,omitempty"`
	SamplingSpaceSizes    googleapi.Int64s        `json:"samplingSpaceSizes,omitempty"`
	Totals                []string                `json:"totals,omitempty"`
}

// ParseReportsResponse turns ga.GetReportsResponse into ParseReport
func ParseReportsResponse(responses []*ga.GetReportsResponse, maxRows int64) (parsedReport *ParseReport, pageToken string) {

	parsed := ParseReport{}
	parsedRowp := make([]*ParseReportRow, maxRows+1)
	rowNum := 0
	fmt.Println("rows to fetch: ", maxRows)

	for _, res := range responses {

		if res.QueryCost > 0 {
			fmt.Println("QueryCost: ", res.QueryCost, " ResourcesQuotasRemaining: ", res.ResourceQuotasRemaining)
		}

		for i, report := range res.Reports {
			//fmt.Println("parse i:", i)
			//js, _ := json.Marshal(report)
			//fmt.Println(string(js))

			if i == 0 {
				var metHeaders []*ga.MetricHeaderEntry
				for _, met := range report.ColumnHeader.MetricHeader.MetricHeaderEntries {
					metHeaders = append(metHeaders, met)
				}

				parsed.ColumnHeaderDimension = report.ColumnHeader.Dimensions
				parsed.ColumnHeaderMetrics = metHeaders
				parsed.DataLastRefreshed = report.Data.DataLastRefreshed
				parsed.IsDataGolden = report.Data.IsDataGolden
				parsed.Maximums = report.Data.Maximums[0].Values
				parsed.Minimums = report.Data.Minimums[0].Values
				parsed.RowCount = report.Data.RowCount
				parsed.SamplesReadCounts = report.Data.SamplesReadCounts
				parsed.SamplingSpaceSizes = report.Data.SamplingSpaceSizes
				parsed.Totals = report.Data.Totals[0].Values
			}

			for _, row := range report.Data.Rows {
				fmt.Println("Parsing row: ", rowNum, row.Dimensions)
				mets := row.Metrics[0].Values
				parsedRowp[rowNum] = &ParseReportRow{Dimensions: row.Dimensions, Metrics: mets}
				rowNum++
			}

			// 0 indexed, only last page of results
			if i == (len(res.Reports) - 1) {
				pageToken = report.NextPageToken
			}
		}
	}

	parsed.Rows = parsedRowp

	// js, _ := json.Marshal(parsed)
	// fmt.Println("parsed: ", string(js))

	return &parsed, pageToken

}
