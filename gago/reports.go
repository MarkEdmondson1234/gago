package gago

import (
	"log"
	"strings"

	ga "google.golang.org/api/analyticsreporting/v4"
)

//makeRequest creates the request(s) for fetchReport
func makeRequest(
	start, end, dimensions, metrics string) *ga.ReportRequest {

	// slice of length 1 of type *ga.DateRange
	daterangep := make([]*ga.DateRange, 1)
	// Fill the 1st element with a pointer to a ga.DateRange
	daterangep[0] = &ga.DateRange{StartDate: start, EndDate: end}

	// a slice of dimension strings
	dimSplit := strings.Split(dimensions, ",")
	// make the slice of length of dimensions
	dimp := make([]*ga.Dimension, len(dimSplit))
	for _, dim := range dimSplit {
		dimp = append(dimp, &ga.Dimension{Name: dim})
	}

	// a slice of dimension strings
	metSplit := strings.Split(metrics, ",")
	// make the slice of length of dimensions
	metp := make([]*ga.Metric, len(metSplit))
	for _, met := range dimSplit {
		metp = append(metp, &ga.Metric{Expression: met})
	}

	requests := ga.ReportRequest{}
	requests.DateRanges = daterangep
	requests.Dimensions = dimp
	requests.Metrics = metp

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
