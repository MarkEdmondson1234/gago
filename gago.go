package gago

import (
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/analytics/v3"
	"google.golang.org/api/analyticsreporting/v4"

	"fmt"
	"io/ioutil"

	csvtag "github.com/artonge/go-csv-tag"
)

// Authenticate Create clients for v4 and v3 Google Analytics API via JSON credentials file
func Authenticate(file string) (*analyticsreporting.Service, *analytics.Service) {
	key, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	jwtConf, err := google.JWTConfigFromJSON(
		key,
		analytics.AnalyticsReadonlyScope,
	)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := jwtConf.Client(oauth2.NoContext)

	//how does this work with NewService?
	analyticsreportingService, err := analyticsreporting.New(httpClient)
	if err != nil {
		log.Fatal(err)
	}
	analyticsService, err := analytics.New(httpClient)
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

// GetAccounts gets the analytics accounts available for this user
func GetAccounts(service *analytics.Service) {

	accountResponse, err := service.Management.Accounts.List().Do()
	if err != nil {
		log.Fatal("Can't find any accounts for this authentication")
	}

	fmt.Println("Found the following accounts:")
	for _, acc := range accountResponse.Items {

		fmt.Println(acc.Id, acc.Name)
	}

}

//GetAccountSummary gets account summary including web properties and viewIds
func GetAccountSummary(service *analytics.Service) {

	accountSummaryResponse, err := service.Management.AccountSummaries.List().Do()
	if err != nil {
		log.Fatal("Can't find account summary")
	}

	fmt.Println("Found following account summary:")
	for _, ass := range accountSummaryResponse.Items {

		fmt.Println("Account:", ass.Id, ass.Name)

		for _, wp := range ass.WebProperties {

			fmt.Println("WebProperty:", wp.Id, wp.Name)

			for _, view := range wp.Profiles {

				fmt.Println("ViewID:", view.Id, view.Name)
			}

		}

	}
}
