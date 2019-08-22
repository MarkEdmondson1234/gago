package gago

import (
	"context"
	"log"

	"google.golang.org/api/option"

	"google.golang.org/api/analytics/v3"
	"google.golang.org/api/analyticsreporting/v4"
)

// Authenticate Create clients for v4 and v3 Google Analytics API via JSON credentials file
// file is a filepath pointing to the location of a service credentials json file downloaded from your GCP console
// Remember to add the service account email to the GA Views you want to download from
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
