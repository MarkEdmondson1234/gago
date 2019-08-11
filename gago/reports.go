package gago

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	ga "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/googleapi"
)

//GoogleAnalytics Make a request to the v4 Reporting API
func GoogleAnalytics(
	service *ga.Service,
	viewID, start, end, dimensions, metrics string,
	maxRows int64,
	useResourceQuotas, antiSample bool) []*ParseReport {

	// init ""
	var pageToken string
	var pageSize, pageLimit, fetchedRows int64

	pageLimit = 10

	pageSize = pageLimit

	if maxRows < pageLimit {
		// if first page needs to fetch less than 10k default
		pageSize = maxRows - 1
	}

	maxPages := (maxRows / pageLimit) + 1
	reqp := make([]*ga.ReportRequest, 5)
	responses := make([]*ga.GetReportsResponse, maxPages)
	parseReportList := make([]*ParseReport, maxPages)

	fmt.Println("maxPages: ", maxPages)

	fetchedRows = 0
	fetchMore := true
	for i := 0; fetchMore; i++ {
		fmt.Println("paging: ", i, fetchMore)

		// a loop around 5 requests
		for j := range reqp {
			fmt.Println("ps", pageSize, " pt", pageToken, " pl", pageLimit, " mr", maxRows, "fr", fetchedRows)
			req := makeRequest(
				viewID,
				start,
				end,
				dimensions,
				metrics,
				pageSize,
				pageToken)
			reqp[j] = req

			pageToken = strconv.FormatInt(fetchedRows+pageSize, 10)
			fetchedRows = fetchedRows + pageSize
			// stop fetching if no pagetoken
			// stop fetching if we adjusted the pageSize down
			// stop fetching if we've reached maxRows
			if pageSize < pageLimit || (maxRows > 0 && fetchedRows >= maxRows) {
				fmt.Println("dont fetchmore")
				fetchMore = false
				break
			}

			// do we need to modify pageSize for next loop?
			if maxRows > 0 && (fetchedRows+pageSize) > maxRows {
				pageSize = maxRows - fetchedRows
			}
		}

		// fetch requests
		// assign this for PArsdReotsResponse dones error
		res1 := fetchReport(service, reqp, useResourceQuotas)
		responses[i] = res1

	}

	for i, res := range responses {
		parsedReport, _ := ParseReportsResponse(res)
		parseReportList[i] = parsedReport
	}

	return parseReportList

}

//makeRequest creates the request(s) for fetchReport
// start and end are YYYY-mm-dd
// dimensions and metrics are ga:dim1,ga:dim2 and ga:metric1,ga:metric2
func makeRequest(
	viewID, start, end, dimensions, metrics string,
	pageSize int64,
	pageToken string) *ga.ReportRequest {

	// slice of length 1 of type *ga.DateRange
	daterangep := make([]*ga.DateRange, 1)
	// Fill the 1st element with a pointer to a ga.DateRange
	daterangep[0] = &ga.DateRange{StartDate: start, EndDate: end}

	dimSplit := strings.Split(dimensions, ",")
	dimp := make([]*ga.Dimension, len(dimSplit))
	for i, dim := range dimSplit {
		dimp[i] = &ga.Dimension{Name: dim}
	}

	metSplit := strings.Split(metrics, ",")
	metp := make([]*ga.Metric, len(metSplit))
	for i, met := range metSplit {
		metp[i] = &ga.Metric{Expression: met}
	}

	// TODO: Make multiple requests based on pagesize
	requests := ga.ReportRequest{}
	requests.DateRanges = daterangep
	requests.Dimensions = dimp
	requests.Metrics = metp
	requests.IncludeEmptyRows = true
	requests.PageSize = pageSize
	requests.ViewId = viewID
	requests.SamplingLevel = "LARGE"
	requests.PageToken = pageToken

	// print out json request
	js, _ := requests.MarshalJSON()
	fmt.Println("Request:", string(js))

	return &requests
}

// FetchReport Perform the GAv4 API request
func fetchReport(
	service *ga.Service,
	reports []*ga.ReportRequest,
	useResourceQuotas bool) *ga.GetReportsResponse {

	reportreq := &ga.GetReportsRequest{ReportRequests: reports, UseResourceQuotas: useResourceQuotas}

	// TODO: parrallise this
	// I don't think this deal with more than 5 reportreq's at once
	report, err := service.Reports.BatchGet(reportreq).Do()
	if err != nil {
		log.Fatal(err)
	}

	js, _ := json.Marshal(report)
	fmt.Println(string(js))

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
func ParseReportsResponse(res *ga.GetReportsResponse) (parsedReport *ParseReport, pageToken string) {

	if res.QueryCost > 0 {
		fmt.Println("QueryCost: ", res.QueryCost, " ResourcesQuotasRemaining: ", res.ResourceQuotasRemaining)
	}

	parsed := ParseReport{}

	for i, report := range res.Reports {

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

		parsedRowp := make([]*ParseReportRow, len(report.Data.Rows))
		for i, row := range report.Data.Rows {
			mets := row.Metrics[0].Values
			parsedRowp[i] = &ParseReportRow{Dimensions: row.Dimensions, Metrics: mets}
		}
		parsed.Rows = parsedRowp

		// 0 indexed, only last page of results
		if i == (len(res.Reports) - 1) {
			pageToken = report.NextPageToken
		}
	}

	return &parsed, pageToken

}
