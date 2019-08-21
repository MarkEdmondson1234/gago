# gago
Google Analytics for Go

Create a CLI that will download GA multi-threaded, using anti-sampling, auto-paging etc. as developed with `googleAnalyticsR::google_analytics()`

Intended use case is for creating executables that can run on any machine without installing another program first, such as R or Python.  This should give more options for running scheduled scripts etc. 

## Command Line Interface

```sh
go install github.com/MarkEdmondson1234/gago/gago
go install github.com/MarkEdmondson1234/gago/gagocli
```

Executable should now be at `$GOPATH/bin/gagocli`

(In future just download the executable once its created)

You can add this to your path variable so you can issue only `gagocli`.

For example on MacOS my $GOPATH/bin is `/Users/me/dev/go/bin`

```sh
sudo nano /etc/paths
# sudo password
# add your $GOPATH/bin to list
```

You can then issue

```sh
gagocli
#gagocli [subcommand...] [arguments...]
#
#Subcommand:
#reports	Download data from Google Analytics API v4
#accounts Get account summary of accounts, webproperties and viewIds
#
#Use -h to get help on subcommand e.g. gagocli report -h
```

Download an auth json file from a GCP project with analytics API enabled, and add the service email to the accounts you want to download.

Supply the auth json file via the `-a` flag or set to `GAGO_AUTH` environment argument in your `~/.bash_profile`

Run via:

```sh
gagocli
#gagocli [subcommand...] [arguments...]
#
#Subcommand:
#reports	Download data from Google Analytics API v4
#accounts Get account summary of accounts, webproperties and viewIds
#
#Use -h to get help on subcommand e.g. gagocli report -h

gagocli reports -h
#Usage of reports:
#  -a string
#    	File path to auth.json service file. Or set via GAGO_AUTH environment argument
#  -antisample
#    	Whether to run anti-sampling
#  -c string
#    	Optional config.yml specifying arguments
#  -dims string
#    	The dimensions ('ga:date,ga:sourceMedium') to run config for
#  -end string
#    	The end date (YYYY-mm-dd) to run config for
#  -max int
#    	The amount of rows to fetch.  Use 0 to fetch all rows (default 1000)
#  -mets string
#    	The metrics ('ga:users,ga:sessions') to run config for
#  -start string
#    	The start date (YYYY-mm-dd) to run config for
#  -view string
#    	The Google Analytics ViewId to run config for

gagocli accounts -h
#Usage of accounts:
#  -a string
#    	File path to auth.json service file. Or set via GAGO_AUTH environment argument
```



## Usage

You can add arguments via the flags of the CLI, or supply a `.yml` file with the configuration of the Google Analytics report to download.  The client email for this file needs to be added to the account/views you want to download as a user.

Example yml file:

```yml
gago:
  view: 81416156
  metrics: ga:sessions,ga:users
  dimensions: ga:date,ga:sourceMedium
  start: 2019-01-01
  end: 2019-08-01
```

This can be sent in the CLI arguments `-c`

```bash
gagocli reports -c config.yml
# {"dimensionHeaderEntries":["ga:date","ga:sourceMedium"],"metricHeaderEntries":[{"name":"ga:sessions","type":"INTEGER"},{"name":"ga:users","type":"INTEGER"}],"values":[{"dimensions":["20190101","(direct) / (none)"]
```

You can override values in the config file via the command line arguments

```bash
gagocli reports -c config.yml -v 123456
# {"dimensionHeaderEntries":["ga:date","ga:sourceMedium"],"metricHeaderEntries":[{"name":"ga:sessions","type":"INTEGER"},{"name":"ga:users","type":"INTEGER"}],"values":[{"dimensions":["20190101","(direct) / (none)"]
```

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

Then get package and dependencies via

```
go get -v github.com/MarkEdmondson1234/gago/gago
go get -v github.com/MarkEdmondson1234/gago/gagocli

go install github.com/MarkEdmondson1234/gago/gago
go install github.com/MarkEdmondson1234/gago/gagocli
```

## Tests

Add the json credential file to an environment argument called `GAGO_AUTH` in your ~/.bash_profile

Then run 

```
go test github.com/MarkEdmondson1234/gago/gago
```

### gago library

The Google Analytics API functions are in a library so it can be used for other Go programs, available via github.com/MarkEdmondson1234/gago/gago

Current functions:

* Authenticate
* GetAccounts
* GetAccountSummary
* GoogleAnalytics
