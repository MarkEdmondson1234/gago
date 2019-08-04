package gago

import (
	"log"

	ga "google.golang.org/api/analyticsreporting/v4"
)

// makeRequest creates the request(s) for fetchReport
// func makeRequest(
// 	start, end string) *ga.GetReportsRequest {

// 	daterangep := ga.DateRange{StartDate: start, EndDate: end}

// 	requests := ga.ReportRequest{}

// 	reportreq := ga.GetReportsRequest{}

// 	return start
// }

// FetchReport Perform the GAv4 API request
func fetchReport(
	service *ga.Service,
	reports *ga.GetReportsRequest) *ga.GetReportsResponse {

	report, err := service.Reports.BatchGet(reports).Do()
	if err != nil {
		log.Fatal(err)
	}

	return report

}
