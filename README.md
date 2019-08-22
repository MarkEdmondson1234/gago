# gago
Google Analytics for Go

[![Build Status](https://travis-ci.com/MarkEdmondson1234/gago.svg?branch=master)](https://travis-ci.com/MarkEdmondson1234/gago)
[![codecov](https://codecov.io/gh/MarkEdmondson1234/gago/branch/master/graph/badge.svg)](https://codecov.io/gh/MarkEdmondson1234/gago)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarkEdmondson1234/gago)](https://goreportcard.com/report/github.com/MarkEdmondson1234/gago)

[gago documentation on godoc.com](https://godoc.org/github.com/MarkEdmondson1234/gago/gago)

## Mission

Create a CLI that will download GA multi-threaded, using anti-sampling, auto-paging etc. as developed with [`googleAnalyticsR::google_analytics()`](https://code.markedmondson.me/googleAnalyticsR/articles/v4.html#anti-sampling)

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
#> gagocli
gagocli [subcommand...] [arguments...]
subcommand:
reports	- Download data from Google Analytics API v4
accounts - Get account summary of accounts, webproperties and viewIds

Use -h to get help on subcommand e.g. gagocli report -h

#> gagocli reports -h
Usage of reports:
  -a string
    	File path to auth.json service file. Or set via GAGO_AUTH environment argument
  -antisample
    	Whether to run anti-sampling
  -c string
    	Optional config.yml specifying arguments
  -dims string
    	The dimensions ('ga:date,ga:sourceMedium') to run config for
  -end string
    	The end date (YYYY-mm-dd) to run config for
  -max int
    	The amount of rows to fetch.  Use 0 to fetch all rows (default 1000)
  -mets string
    	The metrics ('ga:users,ga:sessions') to run config for
  -o string
    	If used will write CSV output to this file
  -start string
    	The start date (YYYY-mm-dd) to run config for
  -v	Verbose output.
  -view string
    	The Google Analytics ViewId to run config for

#> gagocli accounts -h
Usage of accounts:
  -a string
    	File path to auth.json service file. Or set via GAGO_AUTH environment argument
```



## Usage

You can add arguments via the flags of the CLI (see `gagocli reports -h`), or supply a `config.yml` file with the configuration of the Google Analytics report to download.

Example yml file:

```yml
gago:
  view: 1234567
  metrics: ga:sessions,ga:users
  dimensions: ga:date,ga:sourceMedium,ga:landingPagePath
  start: 2019-01-01
  end: 2019-03-01
  maxRows: 1000
  antisample: true
```

This can be sent in the CLI arguments `-c`

```bash
gagocli reports -c config.yml
```

You can override values in the config file via the command line arguments

```bash
gagocli reports -c config.yml -view 1234 -max 10
```

If you have the -v flag (verbose) it writes logs to `stdout`, otherwise it writes to `stderr`, but the output csv format is `stdout`.  This means you can pipe the results to a file:

```bash
gagocli reports -c config.yml > results.csv
cat results.csv
```

Or if you prefer you can specify the file output via the `-o` flag:

```bash
gagocli reports -c config.yml -o results.csv
cat results.csv
```

If you don't want the logs displayed in the console, redirect `stderr` to a file (`2>`) or to null `2>/dev/null`:

```bash
gagocli reports -c config.yml -view 65427188 -max 10 2> mylog.log
cat mylog.log

# no logs
gagocli reports -c config.yml -view 65427188 -max 10 2>/dev/null
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
