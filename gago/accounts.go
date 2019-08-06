package gago

import (
	"log"

	ga "google.golang.org/api/analytics/v3"

	"fmt"
)

// GetAccounts gets the analytics accounts available for this user
func GetAccounts(service *ga.Service) {

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
func GetAccountSummary(service *ga.Service) {

	accountSummaryResponse, err := service.Management.AccountSummaries.List().Do()
	if err != nil {
		log.Fatal("Can't find account summary")
	}

	fmt.Println("Found following account summary:")
	for i, ass := range accountSummaryResponse.Items {
		if i == 0 {
			fmt.Println("accountId accountName webPropertyId webPropertyName viewId viewName")
		}
		for _, wp := range ass.WebProperties {

			for _, view := range wp.Profiles {

				fmt.Println(ass.Id, ass.Name, wp.Id, wp.Name, view.Id, view.Name)
			}

		}

	}
}
