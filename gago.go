package main

import (
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/analytics/v3"
	"google.golang.org/api/analyticsreporting/v4"

	"fmt"
	"io/ioutil"
	"os"

	"github.com/akamensky/argparse"

	"github.com/olebedev/config"

	csvtag "github.com/artonge/go-csv-tag"
)

func authenticate(file string) (*analyticsreporting.Service, *analytics.Service) {
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

type args struct {
	config     string
	auth       string
	view       string
	start      string
	end        string
	antisample bool
}

func parseArgs() args {
	// Create new parser object
	parser := argparse.NewParser("gago", "Downloads data from Google Analytics Reporting API v4")
	// Create flags
	var config = parser.String("c", "config", &argparse.Options{Required: true, Help: "config.yml containing API payload to fetch"})
	var auth = parser.String("a", "auth", &argparse.Options{Required: true, Help: "auth.json service email downloaded from GCP "})
	var view = parser.String("v", "view", &argparse.Options{Required: false, Help: "The Google Analytics ViewId to run config for (Default as configured in config.yml)"})
	var start = parser.String("s", "start", &argparse.Options{Required: false, Help: "The start date (YYYY-mm-dd) to run config for (Default as configured in config.yml)"})
	var end = parser.String("e", "end", &argparse.Options{Required: false, Help: "The end date (YYYY-mm-dd) to run config for (Default as configured in config.yml)"})
	var antisample = parser.Flag("S", "antisample", &argparse.Options{Required: false, Help: "Whether to run anti-sampling (Default as configured in config.yml)"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	return args{
		config:     *config,
		auth:       *auth,
		view:       *view,
		start:      *start,
		end:        *end,
		antisample: *antisample}
}

func readConfigYaml(filename string) *config.Config {

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	cfgg, err := cfg.Get("gago")
	if err != nil {
		panic(err)
	}
	return cfgg
}

func csvOutput(filename string) {
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

func main() {

	args := parseArgs()

	cfg := readConfigYaml(args.config)

	var view = args.view

	if args.view == "" {
		viewid, err := cfg.String("view")
		if err != nil {
			log.Fatal("No viewId passed to fetch data for")
		}
		view = viewid
	} else {
		view = args.view
	}

	t := fmt.Sprintf("Configuration read for viewId: %s", view)
	fmt.Println(t)

	//analyticsreportingService, analyticsService := authenticate(args.auth)
	_, analyticsService := authenticate(args.auth)
	accountResponse, err := analyticsService.Management.Accounts.List().Do()
	if err != nil {
		log.Fatal("Can't find any accounts for this authentication")
	}
	var accountID string

	fmt.Println("Found the following accounts:")
	for i, acc := range accountResponse.Items {

		if i == 0 {
			accountID = acc.Id
		}

		fmt.Println(accountID, acc.Name)
	}

}
