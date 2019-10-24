# gagocli
Google Analytics for Go

## Command Line Interface

Find the latest binary for your system (Windows/MacOS/linux) in the [release page](https://github.com/MarkEdmondson1234/gago/releases)

Download the binary of the appropriate type for your system and put it in your bin folder such as /usr/local/bin - rename it to `gagocli` and chmod to 755

### MacOS

```
curl -o gagocli https://github.com/MarkEdmondson1234/gago/releases/download/v0.2.1/gagocli-v0.2.1-darwin-amd64
# from folder where download is
mv gagocli /usr/local/bin/gagocli
sudo chmod 755 /usr/local/bin/gagocli
```

### Linux

```
sudo apt-get install wget
wget -O gagocli https://github.com/MarkEdmondson1234/gago/releases/download/v0.2.1/gagocli-v0.2.1-linux-amd64
# from folder where download is
mv gagocli /usr/local/bin/gagocli
sudo chmod 755 /usr/local/bin/gagocli
```


### Windows

?
 
## Setup

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
    	The amount of rows to fetch.  Use -1 to fetch all rows (default 1000)
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
