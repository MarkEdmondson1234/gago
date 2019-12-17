package gago

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	ga "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/googleapi"
)

const concurrencyLimit = 10
const apiBatchLimit = 5

var verbose bool

//GoogleAnalyticsRequest Make a request object to pass to GoogleAnalytics
//Service needs to be the reportingService from Analytics
//ViewID is the GA View to fetch from
//Start and End are strings in YYYY-MM-DD format
//Metrics is required, a comma separated string of valid ga: prefix
//Dimensions expects a comma separated string with valid ga: prefix
//MaxRows set to -1 to fetch all rows available
//PageLimit set how many pages to fetch each API request batch - default 10000
//UseResourceQuotas if using GA360, set this to TRUE to get increased quota limits
//AntiSample set to true to attempt to antisample data by breaking up into smaller API calls
//Verbose Prints logs to stdout
type GoogleAnalyticsRequest struct {
	Service                                 *ga.Service
	ViewID, Start, End, Dimensions, Metrics string
	MaxRows, PageLimit                      int64
	UseResourceQuotas, AntiSample, Verbose  bool
	pageSize, fetchedRows, maxPages         int64
	pageToken                               string
	fetchAll                                bool
}

//GoogleAnalytics Make a request to the v4 Reporting API
//Supply the function a GoogleAnalyticsRequest struct object
func GoogleAnalytics(gagoRequest GoogleAnalyticsRequest) *ParseReport {

	if gagoRequest.Verbose {
		verbose = true
		myMessage("verbose=true")
		myMessage(join("GoogleAnalyticsRequest: ", fmt.Sprintf("%+v", gagoRequest)))
	}

	// init ""
	if gagoRequest.PageLimit == 0 {
		gagoRequest.PageLimit = 10000
	}

	if gagoRequest.MaxRows < 0 {
		gagoRequest.MaxRows = gagoRequest.PageLimit
		gagoRequest.fetchAll = true
	}

	if gagoRequest.fetchAll {
		// need to do one fetch to see actual number of rows
		test := GoogleAnalyticsRequest{
			Service:    gagoRequest.Service,
			ViewID:     gagoRequest.ViewID,
			Start:      gagoRequest.Start,
			End:        gagoRequest.End,
			Dimensions: gagoRequest.Dimensions,
			Metrics:    gagoRequest.Metrics,
			MaxRows:    gagoRequest.PageLimit}
		testResponse := GoogleAnalytics(test)

		if gagoRequest.PageLimit > testResponse.RowCount {
			// we are done, return response
			return testResponse
		}

		gagoRequest.MaxRows = testResponse.RowCount
		gagoRequest.fetchAll = true
	}

	gagoRequest.pageSize = gagoRequest.PageLimit

	if gagoRequest.MaxRows < gagoRequest.PageLimit {
		// if first page needs to fetch less than 10k default
		gagoRequest.pageSize = gagoRequest.MaxRows
	}

	var requestList [][]*ga.ReportRequest
	if gagoRequest.AntiSample {
		requestList = makeAntiSampleRequestList(&gagoRequest)
	} else {
		requestList = makeRequestList(&gagoRequest)
	}

	responses := fetchConcurrentReport(requestList, gagoRequest)

	parseReports, _ := parseReportsResponse(responses, gagoRequest)

	return parseReports

}

// ParseReportRow A parsed row of ParseReport
type ParseReportRow struct {
	Dimensions []string `json:"dimensions,omitempty"`
	Metrics    []string `json:"metrics,omitempty"`
}

//ParseReport A parsed Report after all batching and paging
//This takes all the API responses and puts it into a single more usable structure
type ParseReport struct {
	ColumnHeaderDimension []string                `json:"dimensionHeaderEntries,omitempty"`
	ColumnHeaderMetrics   []*ga.MetricHeaderEntry `json:"metricHeaderEntries,omitempty"`
	Rows                  []*ParseReportRow       `json:"values,omitempty"`
	DataLastRefreshed     string                  `json:"dataLastRefreshed,omitempty"`
	IsDataGolden          bool                    `json:"isDataGolden,omitempty"`
	Maximums              []string                `json:"maximums,omitempty"`
	Minimums              []string                `json:"minimums,omitempty"`
	RowCount              int64                   `json:"rowCount,omitempty"`
	FetchedRowCount       int64                   `json:"fetchedRowCount,omitempty"`
	SamplesReadCounts     googleapi.Int64s        `json:"samplesReadCounts,omitempty"`
	SamplingSpaceSizes    googleapi.Int64s        `json:"samplingSpaceSizes,omitempty"`
	Totals                []string                `json:"totals,omitempty"`
}

// parseReportsResponse turns ga.GetReportsResponse into ParseReport
func parseReportsResponse(responses []*ga.GetReportsResponse, gagoRequest GoogleAnalyticsRequest) (parsedReport *ParseReport, pageToken string) {

	parsed := ParseReport{}

	// use append instead as that grows the slice as needed?
	parsedRowp := make([]*ParseReportRow, gagoRequest.PageLimit)
	var rowNum int64

	for _, res := range responses {

		if res == nil {
			//fmt.Println("empty response")
			continue

		}

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
				parsed.SamplesReadCounts = report.Data.SamplesReadCounts
				parsed.SamplingSpaceSizes = report.Data.SamplingSpaceSizes
				parsed.Totals = report.Data.Totals[0].Values
				parsed.RowCount = report.Data.RowCount
			}

			for _, row := range report.Data.Rows {
				//fmt.Println("Parsing row: ", rowNum, row.Dimensions)
				if row == nil {
					continue
				}
				mets := row.Metrics[0].Values
				parsedRowp = append(parsedRowp, &ParseReportRow{Dimensions: row.Dimensions, Metrics: mets})
				rowNum++
			}

			pageToken = report.NextPageToken

		}
	}

	// remove nulls
	parsed.Rows = deleteEmptyRowSlice(parsedRowp)
	parsed.FetchedRowCount = rowNum

	if verbose {
		js, _ := json.Marshal(parsed)
		myMessage(join("Parsed: ", string(js)))
		myMessage(join("Parsed rows:", strconv.FormatInt(rowNum, 10)))
	}

	return &parsed, pageToken

}

//WriteCSV Will write out in CSV format
func WriteCSV(report *ParseReport, file *os.File) {
	// write headers
	var metricHeaders []string
	for _, met := range report.ColumnHeaderMetrics {
		metricHeaders = append(metricHeaders, met.Name)
	}
	headerRow := append(report.ColumnHeaderDimension, metricHeaders...)
	
	file, err := os.Create("goga_data.csv")
    	checkError("Cannot create file.", err)
    	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(headerRow)

	for _, value := range report.Rows {
		// write rows
		fullrow := append(value.Dimensions, value.Metrics...)
		err := writer.Write(fullrow)
		checkError("Couldn't write to file", err)
	}
}
