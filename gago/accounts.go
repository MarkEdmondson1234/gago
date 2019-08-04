package gago

import (
	"log"

	"google.golang.org/api/analytics/v3"

	"fmt"
)

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
