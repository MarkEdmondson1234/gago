package main

import (
	"encoding/csv"
	"flag"

	"github.com/MarkEdmondson1234/gago/gago"

	"fmt"
	"io/ioutil"
	"os"

	"github.com/olebedev/config"
)

type argFlags struct {
	config     string
	auth       string
	view       string
	start      string
	end        string
	metrics    string
	dimensions string
	antisample bool
	maxRows    int64
	output     string
	verbose    bool
}

func parseArgs() (string, argFlags) {

	//report args
	reportCmd := flag.NewFlagSet("reports", flag.ExitOnError)
	var view = reportCmd.String("view", "", "The Google Analytics ViewId to run config for")
	var start = reportCmd.String("start", "", "The start date (YYYY-mm-dd) to run config for")
	var end = reportCmd.String("end", "", "The end date (YYYY-mm-dd) to run config for")
	var antisample = reportCmd.Bool("antisample", false, "Whether to run anti-sampling")
	var metrics = reportCmd.String("mets", "", "The metrics ('ga:users,ga:sessions') to run config for")
	var dimensions = reportCmd.String("dims", "", "The dimensions ('ga:date,ga:sourceMedium') to run config for")
	var maxRows = reportCmd.Int64("max", 1000, "The amount of rows to fetch.  Use 0 to fetch all rows")
	var config = reportCmd.String("c", "", "Optional config.yml specifying arguments")
	var auth = reportCmd.String("a", "", "File path to auth.json service file. Or set via GAGO_AUTH environment argument")
	var output = reportCmd.String("o", "", "If used will write CSV output to this file")
	var verbose = reportCmd.Bool("v", false, "Verbose output.")

	//account args
	accSumCmd := flag.NewFlagSet("accounts", flag.ExitOnError)
	var auth2 = accSumCmd.String("a", "", "File path to auth.json service file. Or set via GAGO_AUTH environment argument")

	if len(os.Args) < 2 {
		usage()
	}

	args := argFlags{}
	switch os.Args[1] {
	case "accounts":
		accSumCmd.Parse(os.Args[2:])
		args.auth = *auth2
		if args.auth == "" {
			if os.Getenv("GAGO_AUTH") == "" {
				fmt.Println("Must supply auth json file via -a or GAGO_AUTH environment arg")
				os.Exit(1)
			}
			args.auth = os.Getenv("GAGO_AUTH")
		}
	case "reports":
		reportCmd.Parse(os.Args[2:])
		args.view = *view
		args.start = *start
		args.end = *end
		args.antisample = *antisample
		args.metrics = *metrics
		args.dimensions = *dimensions
		args.maxRows = *maxRows // default 1000
		args.auth = *auth
		args.config = *config
		args.verbose = *verbose
		args.output = *output

		if args.auth == "" {
			if os.Getenv("GAGO_AUTH") == "" {
				fmt.Println("Must supply auth json file via -a or GAGO_AUTH environment arg")
				os.Exit(1)
			}
			args.auth = os.Getenv("GAGO_AUTH")
		}

		cfg := readConfigYaml(*config)

		if args.view == "" {
			viewid, err := cfg.String("view")
			checkError("No viewId passed to fetch data for", err)
			args.view = viewid
		}

		if args.start == "" {
			start, err := cfg.String("start")
			checkError("No start passed to fetch data for", err)
			args.start = start
		}

		if args.end == "" {
			end, err := cfg.String("end")
			checkError("No end passed to fetch data for", err)
			args.end = end
		}

		if !args.antisample {
			as, _ := cfg.Bool("antisample")
			// if nothing, its ok as default is nothing
			args.antisample = as
		}

		if args.metrics == "" {
			mets, err := cfg.String("metrics")
			checkError("No metrics passed", err)
			args.metrics = mets
		}

		if args.dimensions == "" {
			dims, _ := cfg.String("dimensions")

			args.dimensions = dims
		}

		// will use flag default of 1000 if not in config or flags
		if args.maxRows == 1000 {
			mr, _ := cfg.Int("maxRows")
			if mr != 0 {
				args.maxRows = int64(mr)
			}

		}

	default:
		fmt.Println("Command not recognised:", os.Args[1])
		usage()
	}

	return os.Args[1], args

}

func usage() {
	fmt.Println(usageText)
	os.Exit(0)
}

var usageText = `gagocli [subcommand...] [arguments...]
subcommand:
reports	- Download data from Google Analytics API v4
accounts - Get account summary of accounts, webproperties and viewIds

Use -h to get help on subcommand e.g. gagocli report -h
`

func readConfigYaml(filename string) *config.Config {

	if filename == "" {
		cfg, _ := config.ParseYaml("")
		return cfg
	}
	file, err := ioutil.ReadFile(filename)
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	cfgg, err := cfg.Get("gago")
	if err != nil {
		fmt.Println("Incorrect gago configuration in file:", filename)
		os.Exit(1)
	}

	return cfgg
}

func main() {

	cmd, flags := parseArgs()

	analyticsreportingService, analyticsService := gago.Authenticate(flags.auth)

	switch cmd {
	case "accounts":
		gago.GetAccountSummary(analyticsService)
	case "reports":
		var req = gago.GoogleAnalyticsRequest{
			Service:    analyticsreportingService,
			ViewID:     flags.view,
			Start:      flags.start,
			End:        flags.end,
			Dimensions: flags.dimensions,
			Metrics:    flags.metrics,
			MaxRows:    flags.maxRows,
			AntiSample: flags.antisample,
			Verbose:    flags.verbose}

		report := gago.GoogleAnalytics(req)

		// write headers
		var metricHeaders []string
		for _, met := range report.ColumnHeaderMetrics {
			metricHeaders = append(metricHeaders, met.Name)
		}
		headerRow := append(report.ColumnHeaderDimension, metricHeaders...)

		var writer *csv.Writer
		if len(flags.output) > 0 {
			// write to csv
			file, err := os.Create(flags.output)
			if err != nil {
				fmt.Println("Couldn't open file to write to: ", flags.output)
				os.Exit(1)
			}
			defer file.Close()

			writer = csv.NewWriter(file)
		} else {
			// print to console
			writer = csv.NewWriter(os.Stdout)
		}

		defer writer.Flush()
		writer.Write(headerRow)

		for _, value := range report.Rows {
			// write rows
			fullrow := append(value.Dimensions, value.Metrics...)
			err := writer.Write(fullrow)
			checkError("Couldn't write to file", err)
		}

		//fmt.Println("Downloaded Rows: ", report.RowCount)
	default:
		fmt.Println("Command not recognised:", os.Args[1])
		os.Exit(1)
	}

}

func checkError(s string, err error) {
	if err != nil {
		fmt.Println(s)
		os.Exit(1)
	}
}
