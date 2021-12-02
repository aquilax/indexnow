package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aquilax/indexnow"
)

func main() {
	searchEngineHost := os.Args[1]
	key := os.Args[2]
	urlToAdd := os.Args[3]
	in := indexnow.New(searchEngineHost, &indexnow.Ownership{Key: key}, http.DefaultTransport)
	urlToSubmit := indexnow.GetSingleSubmitUrl(searchEngineHost, key, "", urlToAdd)
	fmt.Println(urlToSubmit)
	_, err := in.SubmitSingleURL(urlToAdd)
	if err != nil {
		panic(err)
	}
}
