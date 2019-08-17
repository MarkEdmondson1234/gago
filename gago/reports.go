package gago

import (
	"fmt"

	ga "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/googleapi"
)

const concurrencyLimit = 10
const apiBatchLimit = 5

//GoogleAnalyticsRequest Make a request object to pass to GoogleAnalytics
type GoogleAnalyticsRequest struct {
	Service                                 *ga.Service
	ViewID, Start, End, Dimensions, Metrics string
	MaxRows, PageLimit                      int64
	UseResourceQuotas, AntiSample           bool
	pageSize, fetchedRows, maxPages         int64
	pageToken                               string
}

//GoogleAnalytics Make a request to the v4 Reporting API
func GoogleAnalytics(gagoRequest GoogleAnalyticsRequest) *ParseReport {

	// init ""
	if gagoRequest.MaxRows == 0 {
		gagoRequest.MaxRows = 1000
	}
	if gagoRequest.PageLimit == 0 {
		gagoRequest.PageLimit = 10000
	}

	gagoRequest.pageSize = gagoRequest.PageLimit
	gagoRequest.MaxRows = gagoRequest.MaxRows - 1 //0 index

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

	//js, _ := json.MarshalIndent(responses, "", " ")
	//fmt.Println("\n# All Responses:", string(js))

	parseReports, _ := parseReportsResponse(responses, gagoRequest.fetchedRows)

	return parseReports

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
func parseReportsResponse(responses []*ga.GetReportsResponse, maxRows int64) (parsedReport *ParseReport, pageToken string) {

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
