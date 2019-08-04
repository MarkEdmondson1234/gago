package gago

import (
	"log"

	ga "google.golang.org/api/analyticsreporting/v4"
)

//makeRequest creates the request(s) for fetchReport
func makeRequest(
	start, end, dimension metric string) *ga.ReportRequest {

	requests := ga.ReportRequest{}
	requests.DateRanges = []*ga.DateRange{{StartDate: start, EndDate: end}}
	requests.Dimensions = []*ga.Dimension{{Name: dimension}}
	requests.Metrics = []*ga.Metric{{Expression: metric}}

	return &requests
}

// FetchReport Perform the GAv4 API request
func fetchReport(
	service *ga.Service,
	reports []*ga.ReportRequest) *ga.GetReportsResponse {

	reportreq := &ga.GetReportsRequest{ReportRequests: reports}

	report, err := service.Reports.BatchGet(reportreq).Do()
	if err != nil {
		log.Fatal(err)
	}

	return report

}
