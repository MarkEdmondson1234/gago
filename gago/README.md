# gago library
Google Analytics for Go

## Development

Get package and dependencies via

```
go get -v github.com/MarkEdmondson1234/gago/gago
go install github.com/MarkEdmondson1234/gago/gago
```

## Tests

Add the json credential file to an environment argument called `GAGO_AUTH` in your ~/.bash_profile

Then run 

```
go test github.com/MarkEdmondson1234/gago/gago
```

### gago library

Current functions:

* Authenticate
* GetAccounts
* GetAccountSummary
* GoogleAnalytics
* WriteCSV