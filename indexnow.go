// package indexnow implements the IndexNow submission protocol as described in
// https://www.indexnow.org/documentation
package indexnow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const MAX_BATCH_SIZE = 10000
const SUBMISSION_CONTENT_TYPE = "application/json; charset=utf-8"

// IndexNow client
type IndexNow struct {
	searchEngineHost string
	key              string
	keyLocation      string
	client           *http.Client
}

// Ownership contains key or key location for claiming resource ownership
type Ownership struct {
	Key         string
	KeyLocation string
}

// SubmissionPayload contains POST payload when submitting multiple URLs
type SubmissionPayload struct {
	Host        string   `json:"host"`
	Key         string   `json:"key"`
	KeyLocation string   `json:"keyLocation,omitempty"`
	URLList     []string `json:"urlList"`
}

// New creates IndexNow client
func New(searchEngineHost string, own *Ownership, rt http.RoundTripper) *IndexNow {
	key := ""
	keyOwnership := ""
	if own != nil {
		key = own.Key
		keyOwnership = own.KeyLocation
	}
	return &IndexNow{
		searchEngineHost,
		key,
		keyOwnership,
		&http.Client{
			Transport: rt,
		},
	}
}

// GetSubmitUrl return URL struct for search engine given host name
func GetSubmitUrl(searchEngineHost string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   searchEngineHost,
		Path:   "indexnow",
	}
}

// GetSingleSubmitUrl returns submission URL for singe URL
// Schema: https://<searchengine>/indexnow?url=url-changed&key=your-key
func GetSingleSubmitUrl(searchEngineHost string, key string, keyLocation string, urlToAdd string) string {

	u := GetSubmitUrl(searchEngineHost)
	q := u.Query()
	q.Set("url", urlToAdd)
	q.Set("key", key)
	if keyLocation != "" {
		q.Set("keyLocation", keyLocation)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// SubmitSingleURL submits single URL for indexing
func (in *IndexNow) SubmitSingleURL(urlToAdd string) (*http.Response, error) {
	urlToSubmit := GetSingleSubmitUrl(in.searchEngineHost, in.key, in.keyLocation, urlToAdd)
	resp, err := in.client.Get(urlToSubmit)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("request returned status code: %d", resp.StatusCode)
	}
	return resp, nil
}

// SubmitBatchURLs submits a batch of URLs for indexing
func (in *IndexNow) SubmitBatchURLs(host string, urlsToAdd []string) (*http.Response, error) {
	if len(urlsToAdd) == 0 {
		return nil, nil
	}
	if len(urlsToAdd) > MAX_BATCH_SIZE {
		return nil, fmt.Errorf("batch size can contain up to %d URLs, %d given", MAX_BATCH_SIZE, len(urlsToAdd))
	}
	u := GetSubmitUrl(in.searchEngineHost)
	b, err := json.Marshal(&SubmissionPayload{
		Host:        host,
		Key:         in.key,
		KeyLocation: in.keyLocation,
		URLList:     urlsToAdd,
	})
	if err != nil {
		return nil, err
	}

	resp, err := in.client.Post(u.String(), SUBMISSION_CONTENT_TYPE, bytes.NewReader(b))
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("request returned status code: %d", resp.StatusCode)
	}
	return resp, nil
}
