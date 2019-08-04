package gago

import (
	"context"
	"log"

	"google.golang.org/api/option"

	"google.golang.org/api/analytics/v3"
	"google.golang.org/api/analyticsreporting/v4"

	csvtag "github.com/artonge/go-csv-tag"
)

// Authenticate Create clients for v4 and v3 Google Analytics API via JSON credentials file
func Authenticate(file string) (*analyticsreporting.Service, *analytics.Service) {

	ctx := context.Background()
	analyticsreportingService, err := analyticsreporting.NewService(ctx, option.WithCredentialsFile(file))
	if err != nil {
		log.Fatal(err)
	}

	analyticsService, err := analytics.NewService(ctx, option.WithCredentialsFile(file))
	if err != nil {
		log.Fatal(err)
	}

	return analyticsreportingService, analyticsService
}

//CsvOutput TODO: Create a CSV output
func CsvOutput(filename string) {
	type Demo struct { // A structure with tags
		Name string  `csv:"name"`
		ID   int     `csv:"ID"`
		Num  float64 `csv:"number"`
	}

	tab := []Demo{ // Create the slice where to put the file content
		Demo{
			Name: "some name",
			ID:   1,
			Num:  42.5,
		},
	}

	err := csvtag.DumpToFile(tab, filename)
	if err != nil {
		log.Fatal("Couldn't write to file")
	}
}
