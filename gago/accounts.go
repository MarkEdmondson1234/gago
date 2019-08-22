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

	for _, acc := range accountResponse.Items {

		fmt.Println(acc.Id, acc.Name)
	}

}

//AccountSummary The account summary return object
type AccountSummary struct {
	accountID       []string
	accountName     []string
	webPropertyID   []string
	webPropertyName []string
	viewID          []string
	viewName        []string
}

//GetAccountSummary Gets account summary including web properties and viewIds
func GetAccountSummary(service *ga.Service) AccountSummary {

	accountSummaryResponse, err := service.Management.AccountSummaries.List().Do()
	if err != nil {
		log.Fatal("Can't find account summary")
	}

	var as AccountSummary

	for i, ass := range accountSummaryResponse.Items {
		if i == 0 {
			fmt.Println("accountId accountName webPropertyId webPropertyName viewId viewName")
		}
		for _, wp := range ass.WebProperties {

			for _, view := range wp.Profiles {

				fmt.Println(ass.Id, ass.Name, wp.Id, wp.Name, view.Id, view.Name)
				as.accountID = append(as.accountID, ass.Id)
				as.accountName = append(as.accountName, ass.Name)
				as.webPropertyID = append(as.webPropertyID, wp.Id)
				as.webPropertyName = append(as.webPropertyName, wp.Name)
				as.viewID = append(as.viewID, view.Id)
				as.viewName = append(as.viewName, view.Name)

			}
		}

	}

	return as
}
