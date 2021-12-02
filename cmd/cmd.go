package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/aquilax/indexnow"
)

func main() {
	var searchEngineHost string
	var key string
	var keyLocation string
	var dryRun bool
	flag.StringVar(&searchEngineHost, "host", "", "Search engine hostname")
	flag.StringVar(&key, "key", "", "Submission key")
	flag.StringVar(&keyLocation, "keyLocation", "", "Submission key location URL")
	flag.BoolVar(&dryRun, "dryRun", false, "Dry run")
	flag.Parse()

	urlToAdd := flag.Arg(0)

	if searchEngineHost == "" || (key == "" && keyLocation == "") || urlToAdd == "" {
		flag.Usage()
		os.Exit(0)
	}

	urlToSubmit := indexnow.GetSingleSubmitUrl(searchEngineHost, key, keyLocation, urlToAdd)
	fmt.Println(urlToSubmit)
	if !dryRun {
		in := indexnow.New(searchEngineHost, &indexnow.Ownership{Key: key, KeyLocation: keyLocation}, http.DefaultTransport)
		_, err := in.SubmitSingleURL(urlToAdd)
		if err != nil {
			panic(err)
		}
	}
}
