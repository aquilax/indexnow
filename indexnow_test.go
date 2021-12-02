// package indexnow implements the IndexNow submission protocol as described in
// https://www.indexnow.org/documentation
package indexnow

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (rtp roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rtp(r)
}

func TestNew(t *testing.T) {
	type args struct {
		searchEngineHost string
		own              *Ownership
		rt               http.RoundTripper
	}
	tests := []struct {
		name string
		args args
		want *IndexNow
	}{
		{
			"initializes the struct correctly",
			args{
				"example.com",
				&Ownership{
					Key:         "key",
					KeyLocation: "keyLocation",
				},
				roundTripFunc(func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader("")),
					}, nil
				}),
			},
			&IndexNow{
				searchEngineHost: "example.com",
				key:              "key",
				keyLocation:      "keyLocation",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.searchEngineHost, tt.args.own, tt.args.rt)
			if got.searchEngineHost != tt.want.searchEngineHost {
				t.Errorf("New().searchHost = %v, want %v", got.searchEngineHost, tt.want.searchEngineHost)
			}
			if got.key != tt.want.key {
				t.Errorf("New().key = %v, want %v", got.key, tt.want.key)
			}
			if got.keyLocation != tt.want.keyLocation {
				t.Errorf("New().keyLocation = %v, want %v", got.keyLocation, tt.want.keyLocation)
			}
		})
	}
}

func TestGetSubmitUrl(t *testing.T) {
	type args struct {
		searchEngineHost string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"generates the correct index url",
			args{
				"www.example.com",
			},
			"https://www.example.com/indexnow",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSubmitUrl(tt.args.searchEngineHost); got.String() != tt.want {
				t.Errorf("GetSubmitUrl() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestGetSingleSubmitUrl(t *testing.T) {
	type args struct {
		searchEngineHost string
		key              string
		keyLocation      string
		urlToAdd         string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"generates correct submit url with key",
			args{
				"www.example.com",
				"aabbcceeff",
				"",
				"https://www.example.org/",
			},
			"https://www.example.com/indexnow?key=aabbcceeff&url=https%3A%2F%2Fwww.example.org%2F",
		},
		{
			"generates correct submit url with keyLocation",
			args{
				"www.example.com",
				"",
				"https://www.example.org/key-location/aabbcceeff.txt",
				"https://www.example.org/",
			},
			"https://www.example.com/indexnow?key=&keyLocation=https%3A%2F%2Fwww.example.org%2Fkey-location%2Faabbcceeff.txt&url=https%3A%2F%2Fwww.example.org%2F",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSingleSubmitUrl(tt.args.searchEngineHost, tt.args.key, tt.args.keyLocation, tt.args.urlToAdd); got != tt.want {
				t.Errorf("GetSingleSubmitUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexNow_SubmitSingleURL(t *testing.T) {
	type fields struct {
		searchEngineHost string
		own              *Ownership
		rt               http.RoundTripper
	}
	type args struct {
		urlToAdd string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
		wantErr        bool
	}{
		{
			"submitting single url works as expected",
			fields{
				"www.example.com",
				&Ownership{Key: "key"},
				roundTripFunc(func(r *http.Request) (*http.Response, error) {
					wantSubmitUrl := "https://www.example.com/indexnow?key=key&url=http%3A%2F%2Fwww.example.org"
					if r.URL.String() != wantSubmitUrl {
						t.Errorf("IndexNow.SubmitSingleURL() url = %v, want %v", r.URL.String(), wantSubmitUrl)
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader("success")),
					}, nil
				}),
			},
			args{"http://www.example.org"},
			200,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := New(
				tt.fields.searchEngineHost,
				tt.fields.own,
				tt.fields.rt,
			)
			got, err := in.SubmitSingleURL(tt.args.urlToAdd)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexNow.SubmitSingleURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.StatusCode, tt.wantStatusCode) {
				t.Errorf("IndexNow.SubmitSingleURL() = %v, want %v", got.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestIndexNow_SubmitBatchURLs(t *testing.T) {
	type fields struct {
		searchEngineHost string
		key              string
		keyLocation      string
		client           *http.Client
	}
	type args struct {
		host      string
		urlsToAdd []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &IndexNow{
				searchEngineHost: tt.fields.searchEngineHost,
				key:              tt.fields.key,
				keyLocation:      tt.fields.keyLocation,
				client:           tt.fields.client,
			}
			got, err := in.SubmitBatchURLs(tt.args.host, tt.args.urlsToAdd)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexNow.SubmitBatchURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexNow.SubmitBatchURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}
