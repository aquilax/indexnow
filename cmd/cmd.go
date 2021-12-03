package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/aquilax/indexnow"
)

func readFile(fileName string) ([]string, error) {
	result := make([]string, 0)
	f, err := os.Open(fileName)

	if err != nil {
		return result, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, nil
}

func main() {
	var searchEngineHost string
	var key string
	var keyLocation string
	var fileName string
	var siteHost string
	var dryRun bool

	flag.StringVar(&searchEngineHost, "host", "", "Search engine hostname")
	flag.StringVar(&key, "key", "", "Submission key")
	flag.StringVar(&keyLocation, "keyLocation", "", "Submission key location URL")

	flag.StringVar(&fileName, "file", "", "File containing multiple url to update (each on new line)")
	flag.StringVar(&siteHost, "siteHost", "", "Submitted URLs hostname")

	flag.BoolVar(&dryRun, "dryRun", false, "Dry run")
	flag.Parse()

	urlToAdd := flag.Arg(0)

	if searchEngineHost == "" || (key == "" && keyLocation == "") || (urlToAdd == "" && fileName == "") {
		flag.Usage()
		os.Exit(0)
	}

	if fileName == "" {
		urlToSubmit := indexnow.GetSingleSubmitUrl(searchEngineHost, key, keyLocation, urlToAdd)
		fmt.Println(urlToSubmit)
	}

	if !dryRun {
		var err error
		var resp *http.Response
		in := indexnow.New(searchEngineHost, &indexnow.Ownership{Key: key, KeyLocation: keyLocation}, http.DefaultTransport)

		if fileName != "" {
			urlsToAdd, err := readFile(fileName)
			if err != nil {
				panic(err)
			}
			resp, err = in.SubmitBatchURLs(siteHost, urlsToAdd)
		} else {
			resp, err = in.SubmitSingleURL(urlToAdd)
		}
		fmt.Printf("response: %#v \n", resp)
		if err != nil {
			panic(err)
		}
	}
}
