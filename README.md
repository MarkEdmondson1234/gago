# gago
Google Analytics for Go

Create a CLI that will download GA multi-threaded, using anti-sampling, auto-paging etc. as developed with `googleAnalyticsR::google_analytics()`

Intended use case is for creating executables that can run on any machine without installing another program first, such as R or Python.  This should give more options for running scheduled scripts etc. 

## References

* https://github.com/avelino/awesome-go#authentication-and-oauth
* https://machiel.me/post/using-google-analytics-api-with-golang/
* https://gobyexample.com/structs
* https://godoc.org/google.golang.org/api/analyticsreporting/v4
* https://godoc.org/golang.org/x/oauth2/google#CredentialsFromJSON
* https://github.com/golang/oauth2/blob/master/google/default.go#L76
* https://golang.org/doc/codewalk/functions/
* [Building for platforms](https://stackoverflow.com/questions/12168873/cross-compile-go-on-osx) https://goreleaser.com/

## Development

Only need to do this if you are working on the code in Go.

See https://golang.org/doc/code.html

From the root go path, say `~/dev/go` as defined in your `~/.bash_profile`

```
# go path
export GOPATH=$HOME/dev/go
```

There is a module and a command line tool (cli)

## Tests

Add the json credential file to an environment argument called `GAGO_AUTH` in your ~/.bash_profile

Then run 

```
go test github.com/MarkEdmondson1234/gago/gago
```

### Command Line Interface

```sh
$> go install github.com/MarkEdmondson1234/gago/gagocli
```

Executable should now be at `~/dev/go/bin/gagocli - run via:

```sh
$> ./bin/gagocli
[-c|--config] is required
usage: gagocli [-h|--help] -c|--config "<value>" -a|--auth "<value>" [-v|--view
            "<value>"] [-s|--start "<value>"] [-e|--end "<value>"]
            [-S|--antisample]

            Downloads data from Google Analytics Reporting API v4

Arguments:

  -h  --help        Print help information
  -c  --config      config.yml containing API payload to fetch
  -a  --auth        auth.json service email downloaded from GCP 
  -v  --view        The Google Analytics ViewId to run config for (Default as
                    configured in config.yml)
  -s  --start       The start date (YYYY-mm-dd) to run config for (Default as
                    configured in config.yml)
  -e  --end         The end date (YYYY-mm-dd) to run config for (Default as
                    configured in config.yml)
  -S  --antisample  Whether to run anti-sampling (Default as configured in
                    config.yml)
```

You can add this to your path variable so you can issue only `gago_cli`

## Use

Requires a `.yml` file with the configuration of the Google Analytics report to download, and a JSON service account credentials file download from GCP console.  The client email for this file needs to be added to the account/views you want to download as a user.

Example yml file:

```yml
gago:
  view: 81416156
  metrics: ga:sessions,ga:users
  dimensions: ga:date,ga:sourceMedium
  start: 2019-01-01
  end: 2019-08-01
```

This and the JSON auth file are required to be sent in the CLI arguments `c` and `a`:

```bash
$> ./bin/gagocli -c config.yml -a your-auth-file.json
Configuration read for viewId: 8141444
Found the following accounts:
474333439 MarkEdmondson
```

### gago library

The Google Analytics API functions are in a library so it can be used for other Go programs, available via github.com/MarkEdmondson1234/gago/gago

Current functions:

* Authenticate
* GetAccounts
* GetAccountSummary
* GoogleAnalytics

Only dimensions and metrics can be fetched for first page at the moment.
